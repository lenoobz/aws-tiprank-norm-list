package main

import (
	"context"
	"log"

	"github.com/google/uuid"
	corid "github.com/hthl85/aws-lambda-corid"
	logger "github.com/hthl85/aws-lambda-logger"
	"github.com/hthl85/aws-tiprank-norm-list/config"
	"github.com/hthl85/aws-tiprank-norm-list/infrastructure/repositories/mongodb/repos"
	"github.com/hthl85/aws-tiprank-norm-list/usecase/assets"
	"github.com/hthl85/aws-tiprank-norm-list/usecase/tiprank"
)

func main() {
	appConf := config.AppConf

	// create new logger
	zap, err := logger.NewZapLogger()
	if err != nil {
		log.Fatal("create app logger failed")
	}
	defer zap.Close()

	// create new repository
	tiprankRepo, err := repos.NewTipRankDividendMongo(nil, zap, &appConf.Mongo)
	if err != nil {
		log.Fatalf("create TipRank dividend mongo failed: %v\n", err)
	}
	defer tiprankRepo.Close()

	// create new repository
	assetRepo, err := repos.NewAssetMongo(nil, zap, &appConf.Mongo)
	if err != nil {
		log.Fatal("create asset mongo failed")
	}
	defer assetRepo.Close()

	// create new service
	tiprankService := tiprank.NewService(tiprankRepo, zap)
	assetService := assets.NewService(assetRepo, *tiprankService, zap)

	// try correlation context
	id, _ := uuid.NewRandom()
	ctx := corid.NewContext(context.Background(), id)

	// try correlation context
	if err := assetService.InsertAssets(ctx); err != nil {
		log.Fatal("insert assets failed")
	}

	// try correlation context
	// if err := assetService.InsertAssetsByTickers(ctx, []string{"TSE:FAP"}); err != nil {
	// 	log.Fatal("insert assets failed")
	// }
}
