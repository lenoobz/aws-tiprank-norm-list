package tiprank

import (
	"context"

	"github.com/lenoobz/aws-tiprank-norm-list/entities"
)

///////////////////////////////////////////////////////////
// Fund Repository Interface
///////////////////////////////////////////////////////////

// Reader interface
type Reader interface {
	FindTipRankAssets(context.Context) ([]*entities.TipRankDividend, error)
	FindTipRankAssetsByTickers(context.Context, []string) ([]*entities.TipRankDividend, error)
}

// Writer interface
type Writer interface{}

// Repo interface
type Repo interface {
	Reader
	Writer
}
