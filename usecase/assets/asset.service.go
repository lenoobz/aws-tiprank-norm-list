package assets

import (
	"context"

	logger "github.com/hthl85/aws-lambda-logger"
	"github.com/hthl85/aws-tiprank-norm-list/usecase/tiprank-assets"
)

// Service sector
type Service struct {
	assetRepo      Repo
	tiprankService tiprank.Service
	log            logger.ContextLog
}

// NewService create new service
func NewService(dividendRepo Repo, tiprankService tiprank.Service, log logger.ContextLog) *Service {
	return &Service{
		assetRepo:      dividendRepo,
		tiprankService: tiprankService,
		log:            log,
	}
}

// InsertAssetDividends adds new assets
func (s *Service) InsertAssets(ctx context.Context, tickers []string) error {
	s.log.Info(ctx, "add new asset")

	assets, err := s.tiprankService.FindTipRankAssets(ctx, tickers)
	if err != nil {
		s.log.Error(ctx, "find all TipRank assets failed", "error", err)
	}

	for _, asset := range assets {
		if err := s.assetRepo.InsertAsset(ctx, asset); err != nil {
			s.log.Error(ctx, "insert asset", "error", err)
			return err
		}
	}

	return nil
}
