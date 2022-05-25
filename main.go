package main

import (
	"bscrap/binance"
	"bscrap/config"
	"bscrap/db"
	"log"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	// https://api.binance.com/api/v3/klines?symbol=ZECUSDT&interval=1w&limit=50&startTime=1621728000000&endTime=1653264000000

	// candleStickDataA := binance.GetCandleStickData("ZECUSDT", "1w", 50, 1621728000000, 1653264000000)
	candleStickDataA, err := binance.GetCandleStickData("ZECUSDT", "1w", 52, 1621728000000, 1653264000000)
	if err != nil {
		log.Panic(err)
	}
	processedDataA := candleStickDataA.ProcessCandleStickData()

	candleStickDataB, err := binance.GetCandleStickData("BTCUSDT", "1w", 52, 1621728000000, 1653264000000)
	if err != nil {
		log.Panic(err)
	}
	processedDataB := candleStickDataB.ProcessCandleStickData()

	rd, err := binance.GetRelation(processedDataA, processedDataB)
	if err != nil {
		log.Panic(err)
	}

	_mongo, err := db.InitMongo(config.DBUri)
	if err != nil {
		log.Panic(err)
	}

	err = _mongo.StoreRelationData(rd)
	if err != nil {
		log.Panic(err)
	}

	interupt := make(chan os.Signal, 1)
	signal.Notify(interupt, syscall.SIGTERM, syscall.SIGINT)
	<-interupt
}
