package entities

import (
	"time"
)

// TipRankAsset struct
type TipRankAsset struct {
	Ticker          string                     `json:"ticker,omitempty"`
	Name            string                     `json:"name,omitempty"`
	Yield           float64                    `json:"yield,omitempty"`
	Amount          float64                    `json:"amount,omitempty"`
	Currency        string                     `json:"currency,omitempty"`
	DividendHistory map[int64]*TipRankDividend `json:"dividendHistory,omitempty"`
}

// TipRankDividend struct
type TipRankDividend struct {
	Dividend       float64    `json:"dividend,omitempty"`
	ExDividendDate *time.Time `json:"exDividendDate,omitempty"`
	RecordDate     *time.Time `json:"recordDate,omitempty"`
	PayoutDate     *time.Time `json:"payoutDate,omitempty"`
}
