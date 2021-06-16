package repos

import (
	"context"
	"fmt"
	"time"

	logger "github.com/hthl85/aws-lambda-logger"
	"github.com/hthl85/aws-tiprank-norm-list/config"
	"github.com/hthl85/aws-tiprank-norm-list/consts"
	"github.com/hthl85/aws-tiprank-norm-list/entities"
	"github.com/hthl85/aws-tiprank-norm-list/infrastructure/repositories/mongodb/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// AssetMongo struct
type AssetMongo struct {
	db     *mongo.Database
	client *mongo.Client
	log    logger.ContextLog
	conf   *config.MongoConfig
}

// NewAssetMongo creates new asset mongo repo
func NewAssetMongo(db *mongo.Database, log logger.ContextLog, conf *config.MongoConfig) (*AssetMongo, error) {
	if db != nil {
		return &AssetMongo{
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

	return &AssetMongo{
		db:     client.Database(conf.Dbname),
		client: client,
		log:    log,
		conf:   conf,
	}, nil
}

// Close disconnect from database
func (r *AssetMongo) Close() {
	ctx := context.Background()
	r.log.Info(ctx, "close mongo client")

	if r.client == nil {
		return
	}

	if err := r.client.Disconnect(ctx); err != nil {
		r.log.Error(ctx, "disconnect mongo failed", "error", err)
	}
}

///////////////////////////////////////////////////////////
// Implement repo interface
///////////////////////////////////////////////////////////

// InsertAsset insert new asset
func (r *AssetMongo) InsertAsset(ctx context.Context, asset *entities.TipRankAsset) error {
	// create new context for the query
	ctx, cancel := createContext(ctx, r.conf.TimeoutMS)
	defer cancel()

	m, err := models.NewAssetModel(ctx, r.log, asset)
	if err != nil {
		r.log.Error(ctx, "create model failed", "error", err)
		return err
	}

	// what collection we are going to use
	colname, ok := r.conf.Colnames[consts.ASSETS_COLLECTION]
	if !ok {
		r.log.Error(ctx, "cannot find collection name")
		return fmt.Errorf("cannot find collection name")
	}
	col := r.db.Collection(colname)

	m.IsActive = true
	m.Schema = r.conf.SchemaVersion
	m.ModifiedAt = time.Now().UTC().Unix()

	filter := bson.D{{
		Key:   "ticker",
		Value: m.Ticker,
	}}

	update := bson.D{
		{
			Key:   "$set",
			Value: m,
		},
		{
			Key: "$setOnInsert",
			Value: bson.D{{
				Key:   "createdAt",
				Value: time.Now().UTC().Unix(),
			}},
		},
	}

	opts := options.Update().SetUpsert(true)

	_, err = col.UpdateOne(ctx, filter, update, opts)
	if err != nil {
		r.log.Error(ctx, "update one failed", "error", err)
		return err
	}

	return nil
}
