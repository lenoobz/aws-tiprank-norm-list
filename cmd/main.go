package main

import (
	"context"
	"log"

	"github.com/google/uuid"
	corid "github.com/hthl85/aws-lambda-corid"
	logger "github.com/hthl85/aws-lambda-logger"
	"github.com/hthl85/aws-tiprank-norm-list/config"
	"github.com/hthl85/aws-tiprank-norm-list/infrastructure/repositories/mongodb/repos"
	"github.com/hthl85/aws-tiprank-norm-list/usecase/tiprank-assets"
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
	tiprankRepo, err := repos.NewTipRankAssetMongo(nil, zap, &appConf.Mongo)
	if err != nil {
		log.Fatalf("create TipRank asset mongo failed: %v\n", err)
	}
	defer tiprankRepo.Close()

	// create new service
	tiprankService := tiprank.NewService(tiprankRepo, zap)

	// try correlation context
	id, _ := uuid.NewRandom()
	ctx := corid.NewContext(context.Background(), id)
	assets, err := tiprankService.FindTipRankAssets(ctx, []string{"TSE:FAP", "TSE:QSR", "TSE:QSP.UN"})
	zap.Info(ctx, "TipRank assets", "assets", assets)
}
