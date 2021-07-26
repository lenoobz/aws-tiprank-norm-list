package assets

import (
	"context"

	"github.com/lenoobz/aws-tiprank-norm-list/entities"
)

///////////////////////////////////////////////////////////
// Stock Repository Interface
///////////////////////////////////////////////////////////

// Reader interface
type Reader interface {
}

// Writer interface
type Writer interface {
	InsertAsset(context.Context, *entities.TipRankDividend) error
}

// Repo interface
type Repo interface {
	Reader
	Writer
}
