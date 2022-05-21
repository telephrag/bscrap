package main

import (
	"bscrap/models"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"math"
	"net/http"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	// url := "https://api.binance.com/api/v3/klines?symbol=ZECUSDT&interval=1h&limit=10&startTime=1653001200000"
	url := "https://api.binance.com/api/v3/klines?symbol=ZECUSDT&interval=1d&limit=30&startTime=1650000000000"

	resp, err := http.Get(url)
	if err != nil {
		log.Fatal(err)
	}
	defer resp.Body.Close()

	content, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatal(err)
	}

	var candleStickData models.CandleStickData
	if err = json.Unmarshal(content, &candleStickData.Data); err != nil {
		log.Fatal(err)
	}

	typicalPriceIntervalData := candleStickData.ProcessCandleStickData()
	fmt.Println(typicalPriceIntervalData.Data)
	fmt.Println(typicalPriceIntervalData.Mean)
	fmt.Println(math.Sqrt(typicalPriceIntervalData.Spread))

	interupt := make(chan os.Signal, 1)
	signal.Notify(interupt, syscall.SIGTERM, syscall.SIGINT)
	<-interupt
}
