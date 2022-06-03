package util

import (
	"bscrap/binance"
	"errors"
	"fmt"
	"strconv"
	"time"
)

func DeduceTime(startTime, endTime, interval, limit string) (int64, int64, error) {
	var st, et int64
	var err error

	i := binance.IntervalLengths[interval]
	r, ok := binance.IntervalRemainders[interval]
	if !ok {
		return -1, -1, errors.New("given interval is invalid or not yet supported for db lookup")
	}

	if startTime != "" {
		if st, err = strconv.ParseInt(startTime, 10, 64); err != nil {
			return -1, -1, fmt.Errorf("%w: \"startTime\" must be int64", err)
		}
		x := i - ((st - r) % i)
		st += x
	}

	if endTime != "" {
		if et, err = strconv.ParseInt(endTime, 10, 64); err != nil {
			return -1, -1, fmt.Errorf("%w: \"endTime\" must be int64", err)
		}
		et = et - (et % i) + r + i - 1
	} else { // assume that user wants data up to current moment...
		et = time.Now().UnixMilli() // ...if endTime and limit are not provided
	}

	if limit != "" {
		l, err := strconv.ParseInt(limit, 10, 64)
		if err != nil {
			return -1, -1, fmt.Errorf("%w: \"limit\" must be int64", err)
		}

		temp := st + i*l - 1 // is off by one error possible here? yep, fixed it
		if temp < et || endTime == "" {
			et = temp
		}
	}

	return st, et, nil
}
