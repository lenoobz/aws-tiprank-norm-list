package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

// TipRankAssetModel struct
type TipRankAssetModel struct {
	ID              *primitive.ObjectID             `bson:"_id,omitempty"`
	CreatedAt       int64                           `bson:"createdAt,omitempty"`
	ModifiedAt      int64                           `bson:"modifiedAt,omitempty"`
	Enabled         bool                            `bson:"enabled"`
	Deleted         bool                            `bson:"deleted"`
	Schema          string                          `bson:"schema,omitempty"`
	Ticker          string                          `bson:"ticker,omitempty"`
	Name            string                          `bson:"name,omitempty"`
	Yield           float64                         `bson:"yield,omitempty"`
	Currency        string                          `bson:"currency,omitempty"`
	DividendHistory map[int64]*DividendHistoryModel `bson:"dividendHistory,omitempty"`
}

// DividendHistoryModel struct
type DividendHistoryModel struct {
	Dividend       float64    `bson:"dividend,omitempty"`
	ExDividendDate *time.Time `bson:"exDividendDate,omitempty"`
	RecordDate     *time.Time `bson:"recordDate,omitempty"`
	DividendDate   *time.Time `bson:"payoutDate,omitempty"`
}
