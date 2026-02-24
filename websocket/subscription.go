package massivews

import (
	"encoding/json"
	"strings"

	"github.com/massive-com/client-go/v3/websocket/models"
	"golang.org/x/exp/maps"
	"golang.org/x/exp/slices"
)

// subscriptions stores topic subscriptions for resubscribing after disconnect.
type subscriptions map[Topic]set
type set map[string]struct{}

// add inserts a set of tickers to the given topic.
func (subs subscriptions) add(topic Topic, tickers ...string) {
	_, prefixExists := subs[topic]
	if !prefixExists || slices.Contains(tickers, "*") {
		subs[topic] = make(set)
	}
	for _, t := range tickers {
		subs[topic][t] = struct{}{}
	}
}

// get retrieves a list of subscription messages based on what has been cached.
func (subs subscriptions) get() []json.RawMessage {
	var msgs []json.RawMessage
	for topic, tickers := range subs {
		msg, err := getSub(models.Subscribe, topic, maps.Keys(tickers)...)
		if err != nil {
			continue // skip malformed messages
		}
		msgs = append(msgs, msg)
	}
	return msgs
}

// delete removes a set of tickers from the given topic.
func (subs subscriptions) delete(topic Topic, tickers ...string) {
	for _, t := range tickers {
		delete(subs[topic], t)
	}
	if len(subs[topic]) == 0 {
		delete(subs, topic)
	}
}

// getSub builds a subscription message for a given topic.
func getSub(action models.Action, topic Topic, tickers ...string) (json.RawMessage, error) {
	if len(tickers) == 0 {
		tickers = []string{"*"}
	}

	var params []string
	for _, ticker := range tickers {
		params = append(params, topic.prefix()+"."+ticker)
	}

	msg, err := json.Marshal(&models.ControlMessage{
		Action: action,
		Params: strings.Join(params, ","),
	})
	if err != nil {
		return nil, err
	}

	return msg, nil
}
