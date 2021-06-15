package main

import (
	"context"
	"log"

	"github.com/google/uuid"
	corid "github.com/hthl85/aws-lambda-corid"
	logger "github.com/hthl85/aws-lambda-logger"
	"github.com/hthl85/aws-tiprank-norm-list/config"
	"github.com/hthl85/aws-tiprank-norm-list/infrastructure/repositories/mongodb/repos"
	tiprankdividends "github.com/hthl85/aws-tiprank-norm-list/usecase/tiprank-dividends"
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
	repo, err := repos.NewTipRankDividendMongo(nil, zap, &appConf.Mongo)
	if err != nil {
		log.Fatalf("create TipRank dividend mongo repo failed: %v\n", err)
	}
	defer repo.Close()

	// create new service
	tiprankService := tiprankdividends.NewService(repo, zap)

	// try correlation context
	id, _ := uuid.NewRandom()
	ctx := corid.NewContext(context.Background(), id)
	dividends, err := tiprankService.FindTipRankDividends(ctx, []string{"TSE:FAP", "TSE:QSR", "TSE:QSP.UN"})
	zap.Info(ctx, "TipRank Dividends", "dividends", dividends)
}
