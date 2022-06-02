package binance

import "go.mongodb.org/mongo-driver/bson"

var IntervalLabels = bson.A{"1m", "3m", "5m", "15m", "30m", "1h", "2h", "4h", "6h", "8h", "12h", "1d", "3d", "1w", "1M"}

func ShorterOrEqualTo(interval interface{}) bson.A {
	for i := range IntervalLabels {
		if interval == IntervalLabels[i] {
			return IntervalLabels[:i+1]
		}
	}
	return bson.A{}
}

const m int64 = 60000
const h int64 = m * 60
const d int64 = h * 24

var IntervalLengths = map[string]int64{
	"1m":  m,
	"3m":  m * 3,
	"5m":  m * 5,
	"15m": m * 15,
	"30m": m * 30,
	"1h":  h,
	"2h":  h * 2,
	"4h":  h * 4,
	"6h":  h * 6,
	"8h":  h * 8,
	"12h": h * 12,
	"1d":  d,
	"3d":  d * 3,
	"1w":  d * 7,
	"1M":  d * 28,
}

var IntervalRemainders = map[string]int64{
	"1d": 0,
	"3d": d * 3,
	"1w": d * 4,
}
