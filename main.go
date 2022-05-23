package main

import (
	"bscrap/models"
	"fmt"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	// url := "https://api.binance.com/api/v3/klines?symbol=ZECUSDT&interval=1d&limit=30&startTime=1650000000000"

	candleStickData := models.GetCandleStickData("ZECUSDT", "1d", 30, 1650000000000, 0)

	typicalPriceIntervalData := candleStickData.ProcessCandleStickData()
	fmt.Println(typicalPriceIntervalData)

	interupt := make(chan os.Signal, 1)
	signal.Notify(interupt, syscall.SIGTERM, syscall.SIGINT)
	<-interupt
}
