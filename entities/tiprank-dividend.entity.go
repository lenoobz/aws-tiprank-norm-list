package entities

import (
	"time"
)

// TipRankDividend struct
type TipRankDividend struct {
	Ticker          string                     `json:"ticker,omitempty"`
	Name            string                     `json:"name,omitempty"`
	Yield           float64                    `json:"yield,omitempty"`
	Amount          float64                    `json:"amount,omitempty"`
	Currency        string                     `json:"currency,omitempty"`
	DividendHistory map[int64]*DividendHistory `json:"dividendHistory,omitempty"`
}

// DividendHistory struct
type DividendHistory struct {
	Dividend       float64    `json:"dividend,omitempty"`
	ExDividendDate *time.Time `json:"exDividendDate,omitempty"`
	RecordDate     *time.Time `json:"recordDate,omitempty"`
	PayoutDate     *time.Time `json:"payoutDate,omitempty"`
}
