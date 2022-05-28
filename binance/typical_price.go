package binance

import (
	"errors"
	"fmt"
	"log"
	"strconv"
)

// MiB =  8388608 bits
// mathematically Count <= 43687.8333... ~ 43687
// binance's limit on amount of records at once = 1000
type TypicalPriceData struct { // 544 + 192 * Count
	Symbol     string
	Data       []typicalPrice
	Mean       float64
	Spread     float64 // s^2
	TradeStart int64   // Data[0].TradeStart
	TradeEnd   int64   // Data[l-1].TradeEnd
	Count      int
}

func (csd *CandleStickData) ProcessCandleStickData() (*TypicalPriceData, error) {
	if len(csd.Data) == 0 {
		return nil, errors.New("no candlestick data to process")
	}

	tpd := &TypicalPriceData{}

	tpd.Symbol = csd.Symbol

	count := float64(len(csd.Data))
	tpd.Data = make([]typicalPrice, len(csd.Data))

	selectiveAverageSquare := 0.0
	for i, interval := range csd.Data {
		processed, err := processCandleStick(&interval)
		if err != nil {
			log.Panic(fmt.Errorf("id.Data[%d]: %w", i, err))
		}

		tpd.Data[i] = *processed

		// slower due to multiple divisions
		// but reduces chances of getting an overflow error
		tpd.Mean += processed.Price / count

		selectiveAverageSquare += processed.Price * processed.Price / count
	}

	// M(x^2) - M^2(x), M() -- expected value
	tpd.Spread = (selectiveAverageSquare - tpd.Mean*tpd.Mean)
	if count > 1 { // idk who would want to calc dispersion for such a small sample size
		ub := count / (count - 1) // unbiasing dispersion
		tpd.Spread *= ub
	}

	tpd.TradeStart = tpd.Data[0].TradeStart
	tpd.TradeEnd = tpd.Data[len(tpd.Data)-1].TradeEnd

	tpd.Count = len(tpd.Data)

	return tpd, nil
}

type typicalPrice struct {
	TradeStart int64
	TradeEnd   int64
	Price      float64
}

func processCandleStick(cs *candleStick) (*typicalPrice, error) {
	tp := 0.0 // (low + high + close) / 3
	temp, err := strconv.ParseFloat(cs.PriceLow, 64)
	if err != nil {
		return nil, fmt.Errorf("i.PriceLow: %w", err)
	}
	tp += temp

	temp, err = strconv.ParseFloat(cs.PriceHigh, 64)
	if err != nil {
		return nil, fmt.Errorf("i.PriceHigh: %w", err)
	}
	tp += temp

	temp, err = strconv.ParseFloat(cs.PriceClose, 64)
	if err != nil {
		return nil, fmt.Errorf("i.PriceClose: %w", err)
	}
	tp += temp
	tp /= 3.0

	return &typicalPrice{
		TradeStart: cs.TradeStart,
		TradeEnd:   cs.TradeEnd,
		Price:      tp,
	}, nil
}
