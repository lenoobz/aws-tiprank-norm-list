package tiprankdividends

import (
	"context"

	"github.com/hthl85/aws-tiprank-norm-list/entities"
)

///////////////////////////////////////////////////////////
// Fund Repository Interface
///////////////////////////////////////////////////////////

// Reader interface
type Reader interface {
	FindTipRankDividends(context.Context, []string) ([]*entities.TipRankDividend, error)
}

// Writer interface
type Writer interface{}

// Repo interface
type Repo interface {
	Reader
	Writer
}
