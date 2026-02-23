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
		Feed:   massivews.LaunchpadFeed,
		Market: massivews.Stocks, // Change Market to match examples below (Stocks, Options, Forex, or Crypto)
		Log:    log,
	})
	if err != nil {
		log.Fatal(err)
	}
	defer c.Close()

	// Stocks - make sure you change "Market: massivews.Stocks" above in the massivews.Config
	err = c.Subscribe(massivews.StocksLaunchpadMinAggs, "*")
	//err = c.Subscribe(massivews.StocksLaunchpadValue, "*")

	// Options - make sure you change "Market: massivews.Options" above in the massivews.Config
	//err = c.Subscribe(massivews.OptionsLaunchpadMinAggs, "O:A230616C00070000")
	//err = c.Subscribe(massivews.OptionsLaunchpadValue, "O:A230616C00070000")

	// Forex - make sure you change "Market: massivews.Forex" above in the massivews.Config
	//err = c.Subscribe(massivews.ForexLaunchpadMinAggs, "*")
	//err = c.Subscribe(massivews.ForexLaunchpadValue, "*")

	// Crypto - make sure you change "Market: massivews.Crypto" above in the massivews.Config
	//err = c.Subscribe(massivews.CryptoLaunchpadMinAggs, "*")
	//err = c.Subscribe(massivews.CryptoLaunchpadValue, "*")

	if err != nil {
		log.Error(err)
		return
	}

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
			case models.CurrencyAgg:
				log.WithFields(logrus.Fields{"currency aggregate": out}).Info()
			case models.EquityAgg:
				log.WithFields(logrus.Fields{"equity aggregate": out}).Info()
			case models.LaunchpadValue:
				log.WithFields(logrus.Fields{"value": out}).Info()
			}
		}
	}
}
