package main

import (
	"bscrap/models"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	// https://api.binance.com/api/v3/klines?symbol=ZECUSDT&interval=1w&limit=50&startTime=1621728000000&endTime=1653264000000

	// candleStickDataA := models.GetCandleStickData("ZECUSDT", "1w", 50, 1621728000000, 1653264000000)
	candleStickDataA, err := models.GetCandleStickData("ZECUSDT", "1w", 50, 1621728000000, 1653264000000)
	if err != nil {
		log.Panic(err)
	}
	processedDataA := candleStickDataA.ProcessCandleStickData()

	candleStickDataB, err := models.GetCandleStickData("ZECUSDT", "1w", 50, 1621728000000, 1653264000000)
	if err != nil {
		log.Panic(err)
	}
	processedDataB := candleStickDataB.ProcessCandleStickData()

	rd, err := models.GetRelation(processedDataA, processedDataB)
	if err != nil {
		log.Panic(err)
	}
	fmt.Printf("%#v", rd)

	interupt := make(chan os.Signal, 1)
	signal.Notify(interupt, syscall.SIGTERM, syscall.SIGINT)
	<-interupt
}
