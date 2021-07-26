package repos

import (
	"context"
	"fmt"
	"time"

	logger "github.com/lenoobz/aws-lambda-logger"
	"github.com/lenoobz/aws-tiprank-norm-list/config"
	"github.com/lenoobz/aws-tiprank-norm-list/consts"
	"github.com/lenoobz/aws-tiprank-norm-list/entities"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// TipRankDividendMongo struct
type TipRankDividendMongo struct {
	db     *mongo.Database
	client *mongo.Client
	log    logger.ContextLog
	conf   *config.MongoConfig
}

// NewTipRankDividendMongo creates new TipRank dividend mongo repo
func NewTipRankDividendMongo(db *mongo.Database, log logger.ContextLog, conf *config.MongoConfig) (*TipRankDividendMongo, error) {
	if db != nil {
		return &TipRankDividendMongo{
			db:   db,
			log:  log,
			conf: conf,
		}, nil
	}

	// set context with timeout from the config
	// create new context for the query
	ctx, cancel := createContext(context.Background(), conf.TimeoutMS)
	defer cancel()

	// set mongo client options
	clientOptions := options.Client()

	// set min pool size
	if conf.MinPoolSize > 0 {
		clientOptions.SetMinPoolSize(conf.MinPoolSize)
	}

	// set max pool size
	if conf.MaxPoolSize > 0 {
		clientOptions.SetMaxPoolSize(conf.MaxPoolSize)
	}

	// set max idle time ms
	if conf.MaxIdleTimeMS > 0 {
		clientOptions.SetMaxConnIdleTime(time.Duration(conf.MaxIdleTimeMS) * time.Millisecond)
	}

	// construct a connection string from mongo config object
	cxnString := fmt.Sprintf("mongodb+srv://%s:%s@%s", conf.Username, conf.Password, conf.Host)

	// create mongo client by making new connection
	client, err := mongo.Connect(ctx, clientOptions.ApplyURI(cxnString))
	if err != nil {
		return nil, err
	}

	return &TipRankDividendMongo{
		db:     client.Database(conf.Dbname),
		client: client,
		log:    log,
		conf:   conf,
	}, nil
}

// Close disconnect from database
func (r *TipRankDividendMongo) Close() {
	ctx := context.Background()
	r.log.Info(ctx, "close mongo client")

	if r.client == nil {
		return
	}

	if err := r.client.Disconnect(ctx); err != nil {
		r.log.Error(ctx, "disconnect mongo failed", "error", err)
	}
}

///////////////////////////////////////////////////////////////////////////////
// Implement interface
///////////////////////////////////////////////////////////////////////////////

// FindTipRankAssetsByTickers finds TipRank assets by tickers
func (r *TipRankDividendMongo) FindTipRankAssetsByTickers(ctx context.Context, tickers []string) ([]*entities.TipRankDividend, error) {
	if len(tickers) < 1 {
		return nil, nil
	}

	uppercaseTickers, err := stringsToUpperCase(tickers)
	if err != nil {
		r.log.Error(ctx, "strings to upper case failed", "error", err)
		return nil, err
	}

	// create new context for the query
	ctx, cancel := createContext(ctx, r.conf.TimeoutMS)
	defer cancel()

	// what collection we are going to use
	colname, ok := r.conf.Colnames[consts.TIPRANK_DIVIDEND_LIST_COLLECTION]
	if !ok {
		r.log.Error(ctx, "cannot find collection name")
		return nil, fmt.Errorf("cannot find collection name")
	}
	col := r.db.Collection(colname)

	// filter
	filter := bson.D{
		{
			Key: "ticker",
			Value: bson.D{{
				Key:   "$in",
				Value: uppercaseTickers,
			}},
		},
	}

	// find options
	findOptions := options.Find()

	cur, err := col.Find(ctx, filter, findOptions)

	// only run defer function when find success
	if cur != nil {
		defer func() {
			if deferErr := cur.Close(ctx); deferErr != nil {
				err = deferErr
			}
		}()
	}

	// find was not succeed
	if err != nil {
		r.log.Error(ctx, "find query failed", "error", err)
		return nil, err
	}

	var tiprankAssets []*entities.TipRankDividend

	// iterate over the cursor to decode document one at a time
	for cur.Next(ctx) {
		// decode cursor to activity model
		var tiprankAsset entities.TipRankDividend
		if err = cur.Decode(&tiprankAsset); err != nil {
			r.log.Error(ctx, "decode failed", "error", err)
			return nil, err
		}

		tiprankAssets = append(tiprankAssets, &tiprankAsset)
	}

	if err := cur.Err(); err != nil {
		r.log.Error(ctx, "iterate over cursor failed", "error", err)
		return nil, err
	}

	return tiprankAssets, nil
}

// FindTipRankAssets finds all TipRank assets
func (r *TipRankDividendMongo) FindTipRankAssets(ctx context.Context) ([]*entities.TipRankDividend, error) {
	// create new context for the query
	ctx, cancel := createContext(ctx, r.conf.TimeoutMS)
	defer cancel()

	// what collection we are going to use
	colname, ok := r.conf.Colnames[consts.TIPRANK_DIVIDEND_LIST_COLLECTION]
	if !ok {
		r.log.Error(ctx, "cannot find collection name")
		return nil, fmt.Errorf("cannot find collection name")
	}
	col := r.db.Collection(colname)

	// filter
	filter := bson.D{}

	// find options
	findOptions := options.Find()

	cur, err := col.Find(ctx, filter, findOptions)

	// only run defer function when find success
	if cur != nil {
		defer func() {
			if deferErr := cur.Close(ctx); deferErr != nil {
				err = deferErr
			}
		}()
	}

	// find was not succeed
	if err != nil {
		r.log.Error(ctx, "find query failed", "error", err)
		return nil, err
	}

	var tiprankAssets []*entities.TipRankDividend

	// iterate over the cursor to decode document one at a time
	for cur.Next(ctx) {
		// decode cursor to activity model
		var tiprankAsset entities.TipRankDividend
		if err = cur.Decode(&tiprankAsset); err != nil {
			r.log.Error(ctx, "decode failed", "error", err)
			return nil, err
		}

		tiprankAssets = append(tiprankAssets, &tiprankAsset)
	}

	if err := cur.Err(); err != nil {
		r.log.Error(ctx, "iterate over cursor failed", "error", err)
		return nil, err
	}

	return tiprankAssets, nil
}
