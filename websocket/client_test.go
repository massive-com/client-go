package massivews

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/gorilla/websocket"
	"github.com/massive-com/client-go/v3/websocket/models"
	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
)

func connect(w http.ResponseWriter, r *http.Request) {
	upgrader := websocket.Upgrader{}
	c, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		return
	}
	defer c.Close()

	for {
		mt, msg, err := c.ReadMessage()
		if err != nil {
			return
		}

		var cm models.ControlMessage
		_ = json.Unmarshal(msg, &cm)
		if cm.Action == "auth" && cm.Params == "good" {
			res := []models.ControlMessage{{EventType: models.EventType{EventType: "status"}, Status: "auth_success"}}
			data, _ := json.Marshal(res)
			err = c.WriteMessage(mt, data)
		} else {
			res := []models.ControlMessage{{EventType: models.EventType{EventType: "status"}, Status: "auth_failed"}}
			data, _ := json.Marshal(res)
			err = c.WriteMessage(mt, data)
		}
		if err != nil {
			return
		}
	}
}

func TestNew(t *testing.T) {
	// successful creation
	c, err := New(Config{
		APIKey: "test",
		Feed:   PolyFeed,
		Market: Options,
	})
	assert.NotNil(t, c)
	assert.Nil(t, err)
	assert.Equal(t, "wss://polyfeed.massive.com/options", c.url)
	assert.Equal(t, &nopLogger{}, c.log)

	// empty config
	c, err = New(Config{})
	assert.Nil(t, c)
	assert.NotNil(t, err)
}

func TestConnectAuthSuccess(t *testing.T) {
	s := httptest.NewServer(http.HandlerFunc(connect))
	defer s.Close()

	log := logrus.New()
	log.SetLevel(logrus.DebugLevel)
	u := "ws" + strings.TrimPrefix(s.URL, "http")
	var retries uint64 = 0
	c, err := New(Config{
		APIKey:     "good",
		Feed:       Feed(u),
		Market:     Market(""),
		Log:        log,
		MaxRetries: &retries,
	})
	assert.NotNil(t, c)
	assert.Nil(t, err)

	defer func() {
		time.Sleep(100 * time.Millisecond)
		c.Close()
	}()

	// closing before connecting shouldn't do anthing
	c.Close()

	// connect successfully
	err = c.Connect()
	assert.Nil(t, err)

	// connecting twice shouldn't do anything
	err = c.Connect()
	assert.Nil(t, err)

	c.Close()
	c.Output()
}

func TestConnectAuthFailure(t *testing.T) {
	s := httptest.NewServer(http.HandlerFunc(connect))
	defer s.Close()

	log := logrus.New()
	log.SetLevel(logrus.DebugLevel)
	u := "ws" + strings.TrimPrefix(s.URL, "http")
	var retries uint64 = 0
	c, err := New(Config{
		APIKey:     "bad",
		Feed:       Feed(u),
		Market:     Market(""),
		Log:        log,
		MaxRetries: &retries,
	})
	assert.NotNil(t, c)
	assert.Nil(t, err)

	defer func() {
		time.Sleep(100 * time.Millisecond)
		c.Close()
	}()

	err = c.Connect()
	assert.Nil(t, err)

	err = <-c.Error()
	assert.NotNil(t, err)
}

func TestConnectRetryFailure(t *testing.T) {
	s := httptest.NewServer(http.HandlerFunc(connect))
	defer s.Close()

	log := logrus.New()
	log.SetLevel(logrus.DebugLevel)
	u := "wss" + strings.TrimPrefix(s.URL, "http") // connecting to wss should fail
	var retries uint64 = 1
	c, err := New(Config{
		APIKey:     "bad",
		Feed:       Feed(u),
		Market:     Market(""),
		Log:        log,
		MaxRetries: &retries,
	})
	assert.NotNil(t, c)
	assert.Nil(t, err)
	err = c.Connect()
	assert.NotNil(t, err)
	c.Close()
}

func TestReconnectCallback(t *testing.T) {
	s := httptest.NewServer(http.HandlerFunc(connect))
	defer s.Close()

	reconnectCallbackCount := 0
	log := logrus.New()
	log.SetLevel(logrus.DebugLevel)
	u := "ws" + strings.TrimPrefix(s.URL, "http")
	var retries uint64 = 0
	c, err := New(Config{
		APIKey:     "good",
		Feed:       Feed(u),
		Market:     Market(""),
		Log:        log,
		MaxRetries: &retries,
		ReconnectCallback: func(err error) {
			assert.Nil(t, err)
			reconnectCallbackCount++
		},
	})
	assert.NotNil(t, c)
	assert.Nil(t, err)
	err = c.Connect()
	assert.Nil(t, err)
	c.reconnect()
	c.Close()
	assert.Equal(t, 1, reconnectCallbackCount)
}

func TestRouteOptionsAdmin(t *testing.T) {
	c, err := New(Config{
		APIKey: "test",
		Feed:   RealTime,
		Market: Options,
	})
	assert.Nil(t, err)

	// frames captured from the live feed, including padded spaces in m that must be preserved
	frames := []string{
		`[{"ev":"AD","i":"1794331356867045949-2270938732","t":1794331356867045949,"q":2270938732,"x":300,"m":"FLEX E   PG    LST 96699 135 P EUR 12.18.26 1122 @ 0.05 2PG"},{"ev":"AD","i":"1794331356867503288-2270938735","t":1794331356867503288,"q":2270938735,"x":300,"m":"FLEX E   PG    LST 97132 135 C EUR 12.18.26 1122 @ 35.25 2PG"}]`,
		`[{"ev":"AD","i":"1794331570129741231-2297917417","t":1794331570129741231,"q":2297917417,"x":323,"m":"FLEX E   AAPL  LST AIaWI 345.000 C EUR 01.15.2027 3000  17.20 2AAPL "},{"ev":"AD","i":"1794331570129773756-2297917418","t":1794331570129773756,"q":2297917418,"x":323,"m":"FLEX E   AAPL  LST AIaWQ 345.000 P EUR 01.15.2027 3000  11.20 2AAPL "}]`,
		// retransmission flag (rt) is only present when true
		`[{"ev":"AD","i":"1794331356867045949-2270938732","t":1794331356867045949,"q":2270938732,"x":300,"m":"FLEX E   PG    LST 96699 135 P EUR 12.18.26 1122 @ 0.05 2PG","rt":true}]`,
	}
	for _, frame := range frames {
		var msgs []json.RawMessage
		err = json.Unmarshal([]byte(frame), &msgs)
		assert.Nil(t, err)
		err = c.route(msgs)
		assert.Nil(t, err)
	}

	expected := []models.OptionsAdmin{
		{
			EventType:      models.EventType{EventType: "AD"},
			ID:             "1794331356867045949-2270938732",
			Exchange:       300,
			Message:        "FLEX E   PG    LST 96699 135 P EUR 12.18.26 1122 @ 0.05 2PG",
			Timestamp:      1794331356867045949,
			SequenceNumber: 2270938732,
		},
		{
			EventType:      models.EventType{EventType: "AD"},
			ID:             "1794331356867503288-2270938735",
			Exchange:       300,
			Message:        "FLEX E   PG    LST 97132 135 C EUR 12.18.26 1122 @ 35.25 2PG",
			Timestamp:      1794331356867503288,
			SequenceNumber: 2270938735,
		},
		{
			EventType:      models.EventType{EventType: "AD"},
			ID:             "1794331570129741231-2297917417",
			Exchange:       323,
			Message:        "FLEX E   AAPL  LST AIaWI 345.000 C EUR 01.15.2027 3000  17.20 2AAPL ",
			Timestamp:      1794331570129741231,
			SequenceNumber: 2297917417,
		},
		{
			EventType:      models.EventType{EventType: "AD"},
			ID:             "1794331570129773756-2297917418",
			Exchange:       323,
			Message:        "FLEX E   AAPL  LST AIaWQ 345.000 P EUR 01.15.2027 3000  11.20 2AAPL ",
			Timestamp:      1794331570129773756,
			SequenceNumber: 2297917418,
		},
		{
			EventType:      models.EventType{EventType: "AD"},
			ID:             "1794331356867045949-2270938732",
			Exchange:       300,
			Message:        "FLEX E   PG    LST 96699 135 P EUR 12.18.26 1122 @ 0.05 2PG",
			Timestamp:      1794331356867045949,
			SequenceNumber: 2270938732,
			Retransmission: true,
		},
	}
	for _, want := range expected {
		out := <-c.Output()
		assert.Equal(t, want, out)
	}
	assert.Equal(t, 0, len(c.Output()))
}
