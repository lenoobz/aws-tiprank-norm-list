package main

import (
	"context"
	"log"

	"github.com/aws/aws-lambda-go/lambda"
	logger "github.com/hthl85/aws-lambda-logger"
	"github.com/hthl85/aws-tiprank-norm-list/config"
	"github.com/hthl85/aws-tiprank-norm-list/infrastructure/repositories/mongodb/repos"
	"github.com/hthl85/aws-tiprank-norm-list/usecase/assets"
	"github.com/hthl85/aws-tiprank-norm-list/usecase/tiprank-assets"
)

type TipRankDividendRequest struct {
	Tickers []string `json:"tickers"`
}

func main() {
	lambda.Start(lambdaHandler)
}

func lambdaHandler(ctx context.Context, req TipRankDividendRequest) {
	log.Println("lambda handler is called")

	appConf := config.AppConf

	// create new logger
	zap, err := logger.NewZapLogger()
	if err != nil {
		log.Fatal("create app logger failed")
	}
	defer zap.Close()

	// create new repository
	tiprankRepo, err := repos.NewTipRankAssetMongo(nil, zap, &appConf.Mongo)
	if err != nil {
		log.Fatalf("create TipRank asset mongo failed: %v\n", err)
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
	if err := assetService.InsertAssets(ctx, req.Tickers); err != nil {
		log.Fatal("insert assets failed")
	}
}
