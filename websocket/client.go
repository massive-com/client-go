package massivews

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/url"
	"strings"
	"sync"
	"time"

	"github.com/cenkalti/backoff/v4"
	"github.com/gorilla/websocket"
	"github.com/massive-com/client-go/v2/websocket/models"
	"golang.org/x/exp/maps"
	"golang.org/x/exp/slices"
	"gopkg.in/tomb.v2"
)

const (
	writeWait      = 5 * time.Second
	pongWait       = 30 * time.Second
	pingPeriod     = pongWait - 5*time.Second // send ping 5 seconds before deadline
	maxMessageSize = 1_000_000                // 1MB
)

// Client defines a client to the Massive WebSocket API.
type Client struct {
	apiKey string
	feed   Feed
	market Market
	url    string

	shouldClose bool
	backoff     backoff.BackOff

	mtx    sync.Mutex
	rwtomb tomb.Tomb
	ptomb  tomb.Tomb

	conn   *websocket.Conn
	rQueue chan json.RawMessage
	wQueue chan json.RawMessage
	subs   subscriptions

	rawData              bool
	bypassRawDataRouting bool
	output               chan any
	err                  chan error

	reconnectCallback func(error)
	log               Logger
}

// New creates a client for the Massive WebSocket API.
func New(config Config) (*Client, error) {
	if err := config.validate(); err != nil {
		return nil, fmt.Errorf("invalid client options: %w", err)
	}

	c := &Client{
		apiKey:               config.APIKey,
		feed:                 config.Feed,
		market:               config.Market,
		backoff:              backoff.NewExponentialBackOff(),
		rQueue:               make(chan json.RawMessage, 10000),
		wQueue:               make(chan json.RawMessage, 1000),
		subs:                 make(subscriptions),
		rawData:              config.RawData,
		bypassRawDataRouting: config.BypassRawDataRouting,
		output:               make(chan any, 100000),
		err:                  make(chan error),
		log:                  config.Log,
		reconnectCallback:    config.ReconnectCallback,
	}

	uri, err := url.Parse(string(c.feed))
	if err != nil {
		return nil, fmt.Errorf("invalid data feed format: %v", err)
	}
	uri.Path = strings.Join([]string{uri.Path, string(c.market)}, "/")
	c.url = uri.String()

	if config.MaxRetries != nil {
		c.backoff = backoff.WithMaxRetries(c.backoff, *config.MaxRetries)
	}

	return c, nil
}

// Connect dials the WebSocket server and starts the read/write and process threads.
// If any subscription messages are pushed before connecting, it will also send those
// to the server.
func (c *Client) Connect() error {
	c.mtx.Lock()
	defer c.mtx.Unlock()

	if c.conn != nil {
		return nil
	}

	notify := func(err error, _ time.Duration) {
		c.log.Errorf(err.Error())
	}
	if err := backoff.RetryNotify(c.connect(false), c.backoff, notify); err != nil {
		return err
	}

	return nil
}

// Subscribe sends a subscription message for a topic and set of tickers. If no
// tickers are passed, it will subscribe to all tickers for a given topic.
func (c *Client) Subscribe(topic Topic, tickers ...string) error {
	c.mtx.Lock()
	defer c.mtx.Unlock()

	if !c.market.supports(topic) {
		return fmt.Errorf("topic '%v' not supported for market '%v'", topic.prefix(), c.market)
	}

	if len(tickers) == 0 || slices.Contains(tickers, "*") {
		tickers = []string{"*"}
	}

	subscribe, err := getSub(models.Subscribe, topic, tickers...)
	if err != nil {
		return err
	}

	c.subs.add(topic, tickers...)
	c.wQueue <- subscribe

	return nil
}

// Unsubscribe sends a message to unsubscribe from a topic and set of tickers. If no
// tickers are passed, it will unsubscribe from all tickers for a given topic.
func (c *Client) Unsubscribe(topic Topic, tickers ...string) error {
	c.mtx.Lock()
	defer c.mtx.Unlock()

	if !c.market.supports(topic) {
		return fmt.Errorf("topic '%v' not supported for market '%v'", topic.prefix(), c.market)
	}

	if len(tickers) == 0 || slices.Contains(tickers, "*") {
		tickers = maps.Keys(c.subs[topic])
	}

	unsubscribe, err := getSub(models.Unsubscribe, topic, tickers...)
	if err != nil {
		return err
	}

	c.subs.delete(topic, tickers...)
	c.wQueue <- unsubscribe

	return nil
}

// Output returns the output queue.
func (c *Client) Output() <-chan any {
	return c.output
}

// Error returns an error channel. If the client hits a fatal error (e.g. auth failed),
// it will push an error to this channel and close the connection.
func (c *Client) Error() <-chan error {
	return c.err
}

// Close attempts to gracefully close the connection to the server.
func (c *Client) Close() {
	c.mtx.Lock()
	defer c.mtx.Unlock()
	c.close(false)
}

func newConn(uri string) (*websocket.Conn, error) {
	conn, res, err := websocket.DefaultDialer.Dial(uri, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to dial server: %w", err)
	} else if res.StatusCode != 101 {
		return nil, errors.New("server failed to switch protocols")
	}

	conn.SetReadLimit(maxMessageSize)
	if err := conn.SetReadDeadline(time.Now().Add(pongWait)); err != nil {
		return nil, fmt.Errorf("failed to set read deadline: %w", err)
	}
	conn.SetPongHandler(func(string) error {
		return conn.SetReadDeadline(time.Now().Add(pongWait))
	})

	return conn, nil
}

func (c *Client) connect(reconnect bool) func() error {
	return func() error {
		// dial the server
		conn, err := newConn(c.url)
		if err != nil {
			return err
		}
		c.conn = conn

		// reset write queue and push auth message
		c.wQueue = make(chan json.RawMessage, 1000)
		auth, err := json.Marshal(models.ControlMessage{
			Action: models.Auth,
			Params: c.apiKey,
		})
		if err != nil {
			return fmt.Errorf("failed to marshal auth message: %w", err)
		}
		c.wQueue <- auth

		// push subscription messages
		subs := c.subs.get()
		for _, msg := range subs {
			c.wQueue <- msg
		}

		// start the threads
		c.rwtomb = tomb.Tomb{}
		c.rwtomb.Go(c.read)
		c.rwtomb.Go(c.write)
		if !reconnect {
			c.ptomb = tomb.Tomb{}
			c.ptomb.Go(c.process)
		}

		return nil
	}
}

func (c *Client) reconnect() {
	c.mtx.Lock()
	defer c.mtx.Unlock()

	if c.shouldClose {
		return
	}

	c.log.Debugf("unexpected disconnect: reconnecting")
	c.close(true)

	notify := func(err error, _ time.Duration) {
		c.log.Errorf(err.Error())
		if c.reconnectCallback != nil {
			c.reconnectCallback(err)
		}
	}
	err := backoff.RetryNotify(c.connect(true), c.backoff, notify)
	if err != nil {
		err = fmt.Errorf("error reconnecting: %w: closing connection", err)
		c.log.Errorf(err.Error())
		c.close(false)
		c.err <- err
	} else {
		// Callback on success.
		if c.reconnectCallback != nil {
			c.reconnectCallback(nil)
		}
	}
}

func (c *Client) closeOutput() {
	close(c.output)
	c.log.Debugf("output channel closed")
}

func (c *Client) close(reconnect bool) {
	if c.conn == nil {
		return
	}

	c.rwtomb.Kill(nil)
	if err := c.rwtomb.Wait(); err != nil {
		c.log.Errorf("r/w threads closed: %v", err)
	}

	if !reconnect {
		c.ptomb.Kill(nil)
		if err := c.ptomb.Wait(); err != nil {
			c.log.Errorf("process thread closed: %v", err)
		}
		c.shouldClose = true
		c.closeOutput()
	}

	if c.conn != nil {
		_ = c.conn.Close()
		c.conn = nil
	}
}

func (c *Client) read() error {
	defer func() {
		c.log.Debugf("read thread closed")
	}()

	for {
		select {
		case <-c.rwtomb.Dying():
			return nil
		default:
			_, msg, err := c.conn.ReadMessage()
			if err != nil {
				if websocket.IsCloseError(err, websocket.CloseNormalClosure) {
					return nil
				} else if websocket.IsUnexpectedCloseError(err, websocket.CloseNormalClosure) {
					return fmt.Errorf("connection closed unexpectedly: %w", err)
				}
				return fmt.Errorf("failed to read message: %w", err)
			}
			if err := c.conn.SetReadDeadline(time.Now().Add(pongWait)); err != nil {
				return fmt.Errorf("failed to set read deadline: %w", err)
			}
			c.rQueue <- msg
		}
	}
}

func (c *Client) write() error {
	ticker := time.NewTicker(pingPeriod)
	defer func() {
		c.log.Debugf("write thread closed")
		ticker.Stop()
		go c.reconnect()
	}()

	for {
		select {
		case <-c.rwtomb.Dying():
			if err := c.conn.WriteControl(websocket.CloseMessage, websocket.FormatCloseMessage(websocket.CloseNormalClosure, ""), time.Now().Add(writeWait)); err != nil {
				return fmt.Errorf("failed to gracefully close: %w", err)
			}
			return nil
		case <-ticker.C:
			if err := c.conn.WriteControl(websocket.PingMessage, nil, time.Now().Add(writeWait)); err != nil {
				return fmt.Errorf("failed to send ping message: %w", err)
			}
		case msg := <-c.wQueue:
			if err := c.conn.SetWriteDeadline(time.Now().Add(writeWait)); err != nil {
				return fmt.Errorf("failed to set write deadline: %w", err)
			}
			if err := c.conn.WriteMessage(websocket.TextMessage, msg); err != nil {
				return fmt.Errorf("failed to send message: %w", err)
			}
		}
	}
}

func (c *Client) process() (err error) {
	defer func() {
		// this client should close if it hits a fatal error (e.g. auth failed)
		c.log.Debugf("process thread closed")
		if err != nil {
			go c.Close()
			c.err <- err
		}
	}()

	for {
		select {
		case <-c.ptomb.Dying():
			return nil
		case data := <-c.rQueue:
			if c.rawData && c.bypassRawDataRouting {
				c.output <- data // push raw bytes to output channel
				continue
			}

			var msgs []json.RawMessage
			if err := json.Unmarshal(data, &msgs); err != nil {
				c.log.Errorf("failed to process raw messages: %v", err)
				continue
			}
			if err := c.route(msgs); err != nil {
				return err
			}
		}
	}
}

func (c *Client) route(msgs []json.RawMessage) error {
	for _, msg := range msgs {
		var ev models.EventType
		err := json.Unmarshal(msg, &ev)
		if err != nil {
			c.log.Errorf("failed to process message: %v", err)
			continue
		}

		switch ev.EventType {
		case "status":
			if err := c.handleStatus(msg); err != nil {
				return err
			}
		default:
			c.handleData(ev.EventType, msg)
		}
	}

	return nil
}

func (c *Client) handleStatus(msg json.RawMessage) error {
	var cm models.ControlMessage
	if err := json.Unmarshal(msg, &cm); err != nil {
		c.log.Errorf("failed to unmarshal message: %v", err)
		return nil
	}

	switch cm.Status {
	case "connected":
		c.log.Debugf("connection successful")
	case "auth_success":
		c.log.Debugf("authentication successful")
	case "auth_failed":
		// this is a fatal error so need to close the connection
		return errors.New("authentication failed: closing connection")
	case "success":
		c.log.Debugf("received a successful status message: %v", sanitize(cm.Message))
	case "error":
		c.log.Errorf("received an error status message: %v", sanitize(cm.Message))
	default:
		c.log.Infof("unknown status message '%v': %v", sanitize(cm.Status), sanitize(cm.Message))
	}

	return nil
}

func (c *Client) handleData(eventType string, msg json.RawMessage) {
	if c.rawData {
		c.output <- msg // push raw JSON to output channel
		return
	}

	switch c.market {
	case Stocks:
		switch eventType {
		case "A", "AM":
			var out models.EquityAgg
			if err := json.Unmarshal(msg, &out); err != nil {
				c.log.Errorf("failed to unmarshal message: %v", err)
				return
			}
			c.output <- out
		case "T":
			var out models.EquityTrade
			if err := json.Unmarshal(msg, &out); err != nil {
				c.log.Errorf("failed to unmarshal message: %v", err)
				return
			}
			c.output <- out
		case "Q":
			var out models.EquityQuote
			if err := json.Unmarshal(msg, &out); err != nil {
				c.log.Errorf("failed to unmarshal message: %v", err)
				return
			}
			c.output <- out
		case "LULD":
			var out models.LimitUpLimitDown
			if err := json.Unmarshal(msg, &out); err != nil {
				c.log.Errorf("failed to unmarshal message: %v", err)
				return
			}
			c.output <- out
		case "NOI":
			var out models.Imbalance
			if err := json.Unmarshal(msg, &out); err != nil {
				c.log.Errorf("failed to unmarshal message: %v", err)
				return
			}
			c.output <- out
		case "FMV":
			var out models.FairMarketValue
			if err := json.Unmarshal(msg, &out); err != nil {
				c.log.Errorf("failed to unmarshal message: %v", err)
				return
			}
			c.output <- out
		case "LV":
			var out models.LaunchpadValue
			if err := json.Unmarshal(msg, &out); err != nil {
				c.log.Errorf("failed to unmarshal message: %v", err)
				return
			}
			c.output <- out
		default:
			c.log.Infof("unknown message type '%s' for market %s", sanitize(eventType), c.market)
		}
	case Options:
		switch eventType {
		case "A", "AM":
			var out models.EquityAgg
			if err := json.Unmarshal(msg, &out); err != nil {
				c.log.Errorf("failed to unmarshal message: %v", err)
				return
			}
			c.output <- out
		case "T":
			var out models.EquityTrade
			if err := json.Unmarshal(msg, &out); err != nil {
				c.log.Errorf("failed to unmarshal message: %v", err)
				return
			}
			c.output <- out
		case "Q":
			var out models.EquityQuote
			if err := json.Unmarshal(msg, &out); err != nil {
				c.log.Errorf("failed to unmarshal message: %v", err)
				return
			}
			c.output <- out
		case "FMV":
			var out models.FairMarketValue
			if err := json.Unmarshal(msg, &out); err != nil {
				c.log.Errorf("failed to unmarshal message: %v", err)
				return
			}
			c.output <- out
		case "LV":
			var out models.LaunchpadValue
			if err := json.Unmarshal(msg, &out); err != nil {
				c.log.Errorf("failed to unmarshal message: %v", err)
				return
			}
			c.output <- out
		default:
			c.log.Infof("unknown message type '%s' for market %s", sanitize(eventType), c.market)
		}
	case Forex:
		switch eventType {
		case "CA", "CAS":
			var out models.CurrencyAgg
			if err := json.Unmarshal(msg, &out); err != nil {
				c.log.Errorf("failed to unmarshal message: %v", err)
				return
			}
			c.output <- out
		case "C":
			var out models.ForexQuote
			if err := json.Unmarshal(msg, &out); err != nil {
				c.log.Errorf("failed to unmarshal message: %v", err)
				return
			}
			c.output <- out
		case "FMV":
			var out models.FairMarketValue
			if err := json.Unmarshal(msg, &out); err != nil {
				c.log.Errorf("failed to unmarshal message: %v", err)
				return
			}
			c.output <- out
		case "LV":
			var out models.LaunchpadValue
			if err := json.Unmarshal(msg, &out); err != nil {
				c.log.Errorf("failed to unmarshal message: %v", err)
				return
			}
			c.output <- out
		default:
			c.log.Infof("unknown message type '%s' for market %s", sanitize(eventType), c.market)
		}
	case Crypto:
		switch eventType {
		case "XA", "XAS":
			var out models.CurrencyAgg
			if err := json.Unmarshal(msg, &out); err != nil {
				c.log.Errorf("failed to unmarshal message: %v", err)
				return
			}
			c.output <- out
		case "XT":
			var out models.CryptoTrade
			if err := json.Unmarshal(msg, &out); err != nil {
				c.log.Errorf("failed to unmarshal message: %v", err)
				return
			}
			c.output <- out
		case "XQ":
			var out models.CryptoQuote
			if err := json.Unmarshal(msg, &out); err != nil {
				c.log.Errorf("failed to unmarshal message: %v", err)
				return
			}
			c.output <- out
		case "XL2":
			var out models.Level2Book
			if err := json.Unmarshal(msg, &out); err != nil {
				c.log.Errorf("failed to unmarshal message: %v", err)
				return
			}
			c.output <- out
		case "FMV":
			var out models.FairMarketValue
			if err := json.Unmarshal(msg, &out); err != nil {
				c.log.Errorf("failed to unmarshal message: %v", err)
				return
			}
			c.output <- out
		case "LV":
			var out models.LaunchpadValue
			if err := json.Unmarshal(msg, &out); err != nil {
				c.log.Errorf("failed to unmarshal message: %v", err)
				return
			}
			c.output <- out
		default:
			c.log.Infof("unknown message type '%s' for market %s", sanitize(eventType), c.market)
		}
	case Indices:
		switch eventType {
		case "A", "AM":
			var out models.EquityAgg
			if err := json.Unmarshal(msg, &out); err != nil {
				c.log.Errorf("failed to unmarshal message: %v", err)
				return
			}
			c.output <- out
		case "V":
			var out models.IndexValue
			if err := json.Unmarshal(msg, &out); err != nil {
				c.log.Errorf("failed to unmarshal message: %v", err)
				return
			}
			c.output <- out
		default:
			c.log.Infof("unknown message type '%s' for market %s", sanitize(eventType), c.market)
		}
	case Futures, FuturesCME, FuturesCBOT, FuturesNYMEX, FuturesCOMEX:
		switch eventType {
		case "A", "AM":
			var out models.FuturesAggregate
			if err := json.Unmarshal(msg, &out); err != nil {
				c.log.Errorf("failed to unmarshal message: %v", err)
				return
			}
			c.output <- out
		case "T":
			var out models.FuturesTrade
			if err := json.Unmarshal(msg, &out); err != nil {
				c.log.Errorf("failed to unmarshal message: %v", err)
				return
			}
			c.output <- out
		case "Q":
			var out models.FuturesQuote
			if err := json.Unmarshal(msg, &out); err != nil {
				c.log.Errorf("failed to unmarshal message: %v", err)
				return
			}
			c.output <- out
		default:
			c.log.Infof("unknown message type '%s' for market %s", sanitize(eventType), c.market)
		}
	default:
		c.log.Infof("unknown market %s", c.market)
	}
}

func sanitize(s string) string {
	return strings.Replace(s, "\n", "", -1)
}
