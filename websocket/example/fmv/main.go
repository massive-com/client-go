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
		Feed:   massivews.BusinessFeed,
		Market: massivews.Stocks,
		Log:    log,
	})
	if err != nil {
		log.Fatal(err)
	}
	defer c.Close()

	// FMV
	_ = c.Subscribe(massivews.BusinessFairMarketValue, "*")

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
			case models.FairMarketValue:
				log.WithFields(logrus.Fields{"fmv": out}).Info()

			default:
				log.WithFields(logrus.Fields{"unknown": out}).Info()
			}
		}
	}
}
