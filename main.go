package main

import (
	"bscrap/models"
	"fmt"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	// https://api.binance.com/api/v3/klines?symbol=ZECUSDT&interval=1w&limit=50&startTime=1621728000000&endTime=1653264000000

	candleStickData := models.GetCandleStickData("ZECUSDT", "1w", 50, 1621728000000, 1653264000000)

	typicalPriceIntervalData := candleStickData.ProcessCandleStickData()
	fmt.Println(typicalPriceIntervalData)

	interupt := make(chan os.Signal, 1)
	signal.Notify(interupt, syscall.SIGTERM, syscall.SIGINT)
	<-interupt
}
