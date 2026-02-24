package main

import (
	"os"
	"os/signal"

	massivews "github.com/massive-com/client-go/v3/websocket"
	"github.com/massive-com/client-go/v3/websocket/models"
	"github.com/sirupsen/logrus"
)

func main() {
	log := logrus.New()
	log.SetLevel(logrus.DebugLevel)
	log.SetFormatter(&logrus.JSONFormatter{})
	c, err := massivews.New(massivews.Config{
		APIKey: os.Getenv("MASSIVE_API_KEY"),
		Feed:   massivews.RealTime,
		Market: massivews.Stocks,
		Log:    log,
	})
	if err != nil {
		log.Fatal(err)
	}
	defer c.Close()

	// aggregates
	//_ = c.Subscribe(massivews.StocksMinAggs, "*")
	//_ = c.Subscribe(massivews.StocksSecAggs, "*")

	// trades
	//_ = c.Subscribe(massivews.StocksTrades, "*")
	_ = c.Subscribe(massivews.StocksTrades, "SPY")

	// quotes
	//_ = c.Subscribe(massivews.StocksQuotes, "*")
	_ = c.Subscribe(massivews.StocksQuotes, "SPY")

	if err := c.Connect(); err != nil {
		log.Error(err)
		return
	}

	sigint := make(chan os.Signal, 1)
	signal.Notify(sigint, os.Interrupt)

	for {
		select {
		case <-sigint:
			return
		case <-c.Error():
			return
		case out, more := <-c.Output():
			if !more {
				return
			}
			switch out.(type) {
			case models.EquityAgg:
				log.WithFields(logrus.Fields{"aggregate": out}).Info()
			case models.EquityTrade:
				log.WithFields(logrus.Fields{"trade": out}).Info()
			case models.EquityQuote:
				log.WithFields(logrus.Fields{"quote": out}).Info()
			}
		}
	}
}
