package binance

import (
	"errors"
)

type CandleStickData struct {
	ID       interface{}
	Symbol   string
	Interval string
	Data     []CandleStick
	FromDB   bool
}

// startTime and endTime must be equal to some values inside Data
func (csd *CandleStickData) Shrink(startTime, endTime int64) (*CandleStickData, error) {

	var left, right int = -1, -1
	for i, cs := range csd.Data {
		if cs.StartTime == startTime {
			left = i
			break
		}
	}

	for i, cs := range csd.Data {
		if cs.EndTime == endTime {
			right = i
		}
	}

	if left == -1 || right == -1 {
		return nil, errors.New("borders of time period didn't match any in given CandleStickData")
	}

	res := &CandleStickData{
		Symbol:   csd.Symbol,
		Interval: csd.Interval,
		Data:     make([]CandleStick, right-left+1),
		FromDB:   csd.FromDB,
	}

	count := copy(res.Data, csd.Data[left:right+1])
	if count == 0 {
		return nil, errors.New("failed to copy data to new shrank CandleStickData")
	}

	return res, nil
}
