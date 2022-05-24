package models

type RelationData struct {
	First       *TypicalPriceData
	Second      *TypicalPriceData
	Covariance  float64
	Correlation float64
}

func GetRelation(first, second *TypicalPriceData) *RelationData {
	return nil
}
