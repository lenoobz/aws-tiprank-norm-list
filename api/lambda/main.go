package main

import (
	"context"
	"log"

	"github.com/aws/aws-lambda-go/lambda"
	logger "github.com/hthl85/aws-lambda-logger"
	"github.com/hthl85/aws-tiprank-norm-list/config"
	"github.com/hthl85/aws-tiprank-norm-list/infrastructure/repositories/mongodb/repos"
	"github.com/hthl85/aws-tiprank-norm-list/usecase/assets"
	tiprankdividends "github.com/hthl85/aws-tiprank-norm-list/usecase/tiprank-dividends"
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
	repo, err := repos.NewTipRankDividendMongo(nil, zap, &appConf.Mongo)
	if err != nil {
		log.Fatalf("create TipRank dividend mongo repo failed: %v\n", err)
	}
	defer repo.Close()

	// create new repository
	assetRepo, err := repos.NewAssetMongo(nil, zap, &appConf.Mongo)
	if err != nil {
		log.Fatal("create asset mongo repo failed")
	}
	defer assetRepo.Close()

	// create new service
	tiprankService := tiprankdividends.NewService(repo, zap)
	assetService := assets.NewService(assetRepo, *tiprankService, zap)

	// try correlation context
	if err := assetService.InsertAssets(ctx, req.Tickers); err != nil {
		log.Fatal("insert asset dividends failed")
	}
}
