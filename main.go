package main

import (
	"bscrap/config"
	"bscrap/db"
	"bscrap/env"
	"bscrap/services"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
)

// https://api.binance.com/api/v3/klines?symbol=ZECUSDT&interval=1w&limit=50&startTime=1621728000000&endTime=1653264000000

func main() {

	mi, err := db.InitMongo(config.DBUri)
	if err != nil {
		log.Panic(err)
		return
	}
	env := &env.Env{Mi: mi}

	go func() {
		err := http.ListenAndServe(
			config.Localhost,
			services.Handle(env),
		)
		if err != nil {
			panic(err)
		}
	}()

	interupt := make(chan os.Signal, 1)
	signal.Notify(interupt, syscall.SIGTERM, syscall.SIGINT)
	<-interupt
}
