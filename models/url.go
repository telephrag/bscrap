package models

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
	return u.argAdd("symbol", symbol)
}

func (u *URI) Interval(interval string) *URI {
	return u.argAdd("interval", interval)
}

func (u *URI) StartTime(startTime uint64) *URI {
	return u.argAdd("startTime", fmt.Sprintf("%d", startTime))
}

func (u *URI) Limit(limit int) *URI {
	if limit <= 0 {
		u.err = errors.New("record limit should be positive")
	}

	return u.argAdd("limit", fmt.Sprintf("%d", limit))
}

func (u *URI) Finalize() string {
	return u.scheme + "://" + u.base + "?" + u.args[:len(u.args)-1]
}
