package util

import (
	"errors"
	"strconv"
)

// URI is tasked to prevent bad requests to API in advance
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

func (u *URI) Limit(limit string) *URI {
	if u.err != nil {
		return u
	}

	if limit != "" {
		limitInt, err := strconv.ParseInt(limit, 10, 64)
		if err != nil {
			u.err = errors.New("\"limit\" must be int")
			return u
		}
		if limitInt < 0 {
			u.err = errors.New("\"limit\" must be positive")
			return u
		}

		return u.argAdd("limit", limit)
	}

	return u
}

func (u *URI) Timeframe(startTime, endTime string) *URI {
	if u.err != nil {
		return u
	}

	var stInt int64
	if startTime != "" {
		stInt, u.err = strconv.ParseInt(startTime, 10, 64)
		if u.err != nil {
			u.err = errors.New("\"startTime\" must be int64")
			return u
		}

		if stInt < 0 {
			u.err = errors.New("\"startTime\" must be positive")
			return u
		}

		u.argAdd("startTime", startTime)
	}

	var etInt int64
	if endTime != "" {
		etInt, u.err = strconv.ParseInt(startTime, 10, 64)
		if u.err != nil {
			u.err = errors.New("\"endTime\" must be int64")
			return u
		}

		if etInt < 0 {
			u.err = errors.New("\"endTime\" must be positive")
			return u
		}

		if stInt > etInt {
			u.err = errors.New("startTime exceeds endTime")
			return u
		}

		return u.argAdd("endTime", endTime)
	}

	return u
}

func (u *URI) String() (string, error) {
	if u.err != nil {
		return "", u.err
	}

	return u.scheme + "://" + u.base + "?" + u.args[:len(u.args)-1], nil
}
