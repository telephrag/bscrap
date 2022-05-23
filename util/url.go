package util

import (
	"errors"
	"fmt"
)

type URI struct {
	scheme string
	base   string
	args   string
	err    error
}

func NewURI(base, scheme string) *URI {
	return &URI{
		scheme: scheme,
		base:   base,
	}
}

func (u *URI) Proceed(dst string) *URI {
	u.base += "/" + dst
	return u
}

func (u *URI) argAdd(name, val string) *URI {
	if val != "" {
		u.args += (name + "=" + val + "&")
	}

	return u
}

func (u *URI) Symbol(symbol string) *URI {
	if u.err != nil {
		return u
	}
	return u.argAdd("symbol", symbol)
}

func contains(arr []string, str string) bool {
	for _, elem := range arr {
		if elem == str {
			return true
		}
	}
	return false
}

func (u *URI) Interval(interval string) *URI {
	if u.err != nil {
		return u
	}

	iSet := []string{"1m", "3m", "5m", "15m", "30m", "1h", "2h", "4h", "6h", "8h", "12h", "1d", "3d", "1w", "1M"}
	if contains(iSet, interval) {
		return u.argAdd("interval", interval)
	}
	u.err = errors.New("incorrect interval provided")
	return u
}

// endTime == 0 is ommited
func (u *URI) Timeframe(startTime, endTime uint64) *URI {
	if u.err != nil {
		return u
	}

	if endTime == 0 {
		return u.argAdd("startTime", fmt.Sprintf("%d", startTime))
	}

	if startTime < endTime {
		return u.argAdd("startTime", fmt.Sprintf("%d", startTime)).argAdd("endTime", fmt.Sprintf("%d", endTime))
	}

	u.err = errors.New("startTime exceeds endTime")
	return u
}

func (u *URI) Limit(limit int) *URI {
	if u.err != nil {
		return u
	}

	if limit <= 0 {
		u.err = errors.New("record limit should be positive")
	}

	return u.argAdd("limit", fmt.Sprintf("%d", limit))
}

func (u *URI) String() (string, error) {
	if u.err != nil {
		return "", u.err
	}

	return u.scheme + "://" + u.base + "?" + u.args[:len(u.args)-1], nil
}
