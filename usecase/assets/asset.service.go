package assets

import (
	"context"

	logger "github.com/hthl85/aws-lambda-logger"
	tiprankdividends "github.com/hthl85/aws-tiprank-norm-list/usecase/tiprank-dividends"
)

// Service sector
type Service struct {
	dividendRepo   Repo
	tiprankService tiprankdividends.Service
	log            logger.ContextLog
}

// NewService create new service
func NewService(dividendRepo Repo, tiprankService tiprankdividends.Service, log logger.ContextLog) *Service {
	return &Service{
		dividendRepo:   dividendRepo,
		tiprankService: tiprankService,
		log:            log,
	}
}

// InsertAssetDividends adds new assets
func (s *Service) InsertAssets(ctx context.Context, tickers []string) error {
	s.log.Info(ctx, "add new asset")

	dividends, err := s.tiprankService.FindTipRankDividends(ctx, tickers)
	if err != nil {
		s.log.Error(ctx, "find all TipRank dividends failed", "error", err)
	}

	for _, dividend := range dividends {
		if err := s.dividendRepo.InsertAsset(ctx, dividend); err != nil {
			s.log.Error(ctx, "insert asset", "error", err)
			return err
		}
	}

	return nil
}
