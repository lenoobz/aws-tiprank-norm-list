package assets

import (
	"context"

	logger "github.com/hthl85/aws-lambda-logger"
	"github.com/hthl85/aws-tiprank-norm-list/usecase/tiprank"
)

// Service sector
type Service struct {
	assetRepo      Repo
	tiprankService tiprank.Service
	log            logger.ContextLog
}

// NewService create new service
func NewService(assetRepo Repo, tiprankService tiprank.Service, log logger.ContextLog) *Service {
	return &Service{
		assetRepo:      assetRepo,
		tiprankService: tiprankService,
		log:            log,
	}
}

// InsertAssetDividends adds new assets by tickers
func (s *Service) InsertAssetsByTickers(ctx context.Context, tickers []string) error {
	s.log.Info(ctx, "adding new asset", "tickers", tickers)

	dividends, err := s.tiprankService.FindTipRankDividendsByTickers(ctx, tickers)
	if err != nil {
		s.log.Error(ctx, "find all TipRank dividends by tickers failed", "error", err)
	}

	for _, dividend := range dividends {
		if err := s.assetRepo.InsertAsset(ctx, dividend); err != nil {
			s.log.Error(ctx, "insert asset failed", "error", err)
			return err
		}
	}

	return nil
}

// InsertAssets adds new assets
func (s *Service) InsertAssets(ctx context.Context) error {
	s.log.Info(ctx, "adding new asset")

	dividends, err := s.tiprankService.FindTipRankDividends(ctx)
	if err != nil {
		s.log.Error(ctx, "find all TipRank dividends failed", "error", err)
	}

	for _, dividend := range dividends {
		if err := s.assetRepo.InsertAsset(ctx, dividend); err != nil {
			s.log.Error(ctx, "insert asset failed", "error", err)
			return err
		}
	}

	return nil
}
