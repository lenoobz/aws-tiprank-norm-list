package entities

import (
	"context"
	"time"

	logger "github.com/hthl85/aws-lambda-logger"
	"github.com/hthl85/aws-tiprank-norm-list/utils/ticker"
)

// TipRankDividend struct
type TipRankDividend struct {
	Ticker          string                   `json:"ticker,omitempty"`
	Name            string                   `json:"name,omitempty"`
	Yield           float64                  `json:"yield,omitempty"`
	Amount          float64                  `json:"amount,omitempty"`
	Currency        string                   `json:"currency,omitempty"`
	DividendHistory map[int64]*DividendModel `json:"dividendHistory,omitempty"`
}

// DividendModel struct
type DividendModel struct {
	Dividend       float64    `json:"dividend,omitempty"`
	ExDividendDate *time.Time `json:"exDividendDate,omitempty"`
	RecordDate     *time.Time `json:"recordDate,omitempty"`
	DividendDate   *time.Time `json:"payoutDate,omitempty"`
}

// MapTipRankDividendToAsset map TipRank dividend to asset
func (f *TipRankDividend) MapTipRankDividendToAsset(ctx context.Context, log logger.ContextLog) *Asset {
	asset := &Asset{
		Ticker:           ticker.GenYahooTickerFromTipRankTicker(f.Ticker),
		Name:             f.Name,
		AllocationStock:  0,
		AllocationBond:   0,
		AllocationCash:   0,
		DividendSchedule: "",
		Yield12Month:     0,
		DistYield:        0,
		DistAmount:       0,
	}

	return asset
}
