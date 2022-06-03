package binance

import (
	"errors"
)

type KLineData struct {
	ID       interface{} `json:"id"`
	Symbol   string      `json:"symbol"`
	Interval string      `json:"interval"`
	Data     []KLine     `json:"data"`
	FromDB   bool        `json:"fromDB,omitempty"`
}

// startTime and endTime must be equal to some values inside Data
func (kld *KLineData) Shrink(startTime, endTime int64) (*KLineData, error) {

	var left, right int = -1, -1
	for i, cs := range kld.Data {
		if cs.StartTime == startTime {
			left = i
			break
		}
	}

	for i, cs := range kld.Data {
		if cs.EndTime == endTime {
			right = i
		}
	}

	if left == -1 || right == -1 {
		return nil, errors.New("borders of time period didn't match any in given CandleStickData")
	}

	res := &KLineData{
		ID:       kld.ID,
		Symbol:   kld.Symbol,
		Interval: kld.Interval,
		Data:     make([]KLine, right-left+1),
		FromDB:   kld.FromDB,
	}

	count := copy(res.Data, kld.Data[left:right+1])
	if count == 0 {
		return nil, errors.New("failed to copy data to new shrank CandleStickData")
	}

	return res, nil
}

// mutates kld instead of creating hard copy
func (kld *KLineData) ShrinkSelf(startTime, endTime int64) (*KLineData, error) {

	var left, right int = -1, -1
	for i, cs := range kld.Data {
		if cs.StartTime == startTime {
			left = i
			break
		}
	}

	for i, cs := range kld.Data {
		if cs.EndTime == endTime {
			right = i
		}
	}

	if left == -1 || right == -1 {
		return nil, errors.New("borders of time period didn't match any in given CandleStickData")
	}

	kld.Data = kld.Data[left : right+1]

	return kld, nil
}
