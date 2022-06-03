package main

import (
	"bscrap/bscrap_srv"
	"bscrap/config"
	"bscrap/db"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
)

// https://api.binance.com/api/v3/klines?symbol=ZECUSDT&interval=1w&limit=52&startTime=1621728000000&endTime=1653264000000

func main() {

	mi, err := db.ConnectMongo(config.DBUri)
	if err != nil {
		log.Panic(err)
		return
	}

	bScrapEnv := &bscrap_srv.Env{Mi: mi}
	go func() {
		err := http.ListenAndServe(
			"localhost:8080",
			bscrap_srv.Run(bScrapEnv),
		)
		if err != nil {
			log.Panic(err)
		}
	}()

	interupt := make(chan os.Signal, 1)
	signal.Notify(interupt, syscall.SIGTERM, syscall.SIGINT)
	<-interupt
}
