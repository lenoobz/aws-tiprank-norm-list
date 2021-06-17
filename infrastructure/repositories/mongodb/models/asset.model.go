package models

import (
	"context"
	"sort"
	"strings"
	"time"

	logger "github.com/hthl85/aws-lambda-logger"
	"github.com/hthl85/aws-tiprank-norm-list/consts"
	"github.com/hthl85/aws-tiprank-norm-list/entities"
	"github.com/hthl85/aws-tiprank-norm-list/utils/ticker"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// AssetModel struct
type AssetModel struct {
	ID               *primitive.ObjectID `bson:"_id,omitempty"`
	CreatedAt        int64               `bson:"createdAt,omitempty"`
	ModifiedAt       int64               `bson:"modifiedAt,omitempty"`
	Enabled          bool                `bson:"enabled"`
	Deleted          bool                `bson:"deleted"`
	Schema           string              `bson:"schema,omitempty"`
	Source           string              `bson:"source,omitempty"`
	Ticker           string              `bson:"ticker,omitempty"`
	Name             string              `bson:"name,omitempty"`
	Type             string              `bson:"type,omitempty"`
	AssetClass       string              `bson:"assetClass,omitempty"`
	Currency         string              `bson:"currency,omitempty"`
	AllocationStock  float64             `bson:"allocationStock,omitempty"`
	AllocationBond   float64             `bson:"allocationBond,omitempty"`
	AllocationCash   float64             `bson:"allocationCash,omitempty"`
	DividendSchedule string              `bson:"dividendSchedule,omitempty"`
	Yield12Month     float64             `bson:"yield12Month,omitempty"`
	DistYield        float64             `bson:"distYield,omitempty"`
	DistAmount       float64             `bson:"distAmount,omitempty"`
}

// NewAssetModel create stock model
func NewAssetModel(ctx context.Context, l logger.ContextLog, e *entities.TipRankDividend, schemaVersion string) (*AssetModel, error) {
	var m = &AssetModel{
		ModifiedAt: time.Now().UTC().Unix(),
		Enabled:    true,
		Deleted:    false,
		Schema:     schemaVersion,
	}

	m.Source = strings.ToUpper(consts.DATA_SOURCE)
	m.Type = strings.ToUpper(consts.SECURITY_TYPE)
	m.AssetClass = strings.ToUpper(consts.ASSET_CLASS)

	m.Ticker = ticker.GenYahooTickerFromTipRankTicker(e.Ticker)

	if e.Name != "" {
		m.Name = e.Name
	}

	if e.Currency != "" {
		m.Currency = strings.ToUpper(e.Currency)
	}

	m.AllocationStock = 100
	m.AllocationBond = 0
	m.AllocationCash = 0

	// m.Yield12Month = e.Yield12Month
	m.DistYield = e.Yield
	m.DistAmount = e.Amount
	m.DividendSchedule = calculateDividendSchedule(e.DividendHistory)

	return m, nil
}

func calculateDividendSchedule(dividends map[int64]*entities.DividendHistory) string {
	type kv struct {
		Key   int64
		Value time.Time
	}

	var ss []kv
	for k, v := range dividends {
		ss = append(ss, kv{k, *v.ExDividendDate})
	}

	if len(ss) < 2 {
		return "UNKNOWN"
	}

	sort.Slice(ss, func(i, j int) bool {
		return ss[i].Value.Sub(ss[j].Value) > 0
	})

	diff := ss[0].Value.Sub(ss[1].Value).Seconds() / 86400

	if diff > 80 && diff < 100 {
		return "QUARTERLY"
	}

	if diff > 25 && diff < 35 {
		return "MONTHLY"
	}

	if diff > 175 && diff > 195 {
		return "BI-ANNUALLY"
	}

	if diff > 355 && diff < 375 {
		return "ANNUALLY"
	}

	return "UNKNOWN"
}
