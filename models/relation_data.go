package models

import (
	"errors"
	"math"
)

type RelationData struct {
	First       *TypicalPriceData
	Second      *TypicalPriceData
	Covariance  float64
	Correlation float64
}

func GetRelation(a, b *TypicalPriceData) (*RelationData, error) {

	if a.Count != b.Count {
		return nil, errors.New("dataset sizes mismatch")
	}

	if a.Data[0].TradeStart != b.Data[0].TradeStart ||
		a.Data[0].TradeEnd != b.Data[0].TradeEnd {
		return nil, errors.New("interval sizes mismatch")
	}

	var cov float64
	for i := range a.Data {
		cov += (a.Data[i].Price - a.Mean) * (b.Data[i].Price - b.Mean) / float64(a.Count)
	}

	cor := cov / (math.Sqrt(a.Spread) * math.Sqrt(b.Spread))

	return &RelationData{
		First:       a,
		Second:      b,
		Covariance:  cov,
		Correlation: cor,
	}, nil
}
