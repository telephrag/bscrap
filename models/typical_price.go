package models

import (
	"fmt"
	"log"
	"strconv"
)

type TypicalPrice struct {
	TradeStart uint64
	TradeEnd   uint64
	Price      float64
}

func ProcessCandleStick(cs *CandleStick) (*TypicalPrice, error) {
	typicalPrice := 0.0 // (low + high + close) / 3
	temp, err := strconv.ParseFloat(cs.PriceLow, 64)
	if err != nil {
		return nil, fmt.Errorf("i.PriceLow: %w", err)
	}
	typicalPrice += temp

	temp, err = strconv.ParseFloat(cs.PriceHigh, 64)
	if err != nil {
		return nil, fmt.Errorf("i.PriceHigh: %w", err)
	}
	typicalPrice += temp

	temp, err = strconv.ParseFloat(cs.PriceClose, 64)
	if err != nil {
		return nil, fmt.Errorf("i.PriceClose: %w", err)
	}
	typicalPrice += temp
	typicalPrice /= 3.0

	return &TypicalPrice{
		TradeStart: cs.TradeStart,
		TradeEnd:   cs.TradeEnd,
		Price:      typicalPrice,
	}, nil
}

type TypicalPriceData struct {
	Symbol     string
	Data       []TypicalPrice
	Mean       float64
	Spread     float64 // s^2
	TradeStart uint64  // Data[0].TradeStart
	TradeEnd   uint64  // Data[l-1].TradeEnd
}

func (csd *CandleStickData) ProcessCandleStickData() *TypicalPriceData {
	tpd := &TypicalPriceData{}

	tpd.Symbol = csd.Symbol

	count := float64(len(csd.Data))
	tpd.Data = make([]TypicalPrice, len(csd.Data))

	selectiveAverageSquare := 0.0
	for i, interval := range csd.Data {
		processed, err := ProcessCandleStick(&interval)
		if err != nil {
			log.Panic(fmt.Errorf("id.Data[%d]: %w", i, err))
		}

		tpd.Data[i] = *processed

		// slower due to multiple divisions
		// but reduces chances of overflow error
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

	return tpd
}
