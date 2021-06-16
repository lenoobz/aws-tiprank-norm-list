package tiprank

import (
	"context"

	logger "github.com/hthl85/aws-lambda-logger"
	"github.com/hthl85/aws-tiprank-norm-list/entities"
)

// Service sector
type Service struct {
	repo Repo
	log  logger.ContextLog
}

// NewService create new service
func NewService(repo Repo, log logger.ContextLog) *Service {
	return &Service{
		repo: repo,
		log:  log,
	}
}

// FindTipRankAssetsByTickers finds all tip rank assets by given tickers
func (s *Service) FindTipRankAssetsByTickers(ctx context.Context, tickers []string) ([]*entities.TipRankAsset, error) {
	s.log.Info(ctx, "find TipRank assets by tickers", "tickers", tickers)
	return s.repo.FindTipRankAssetsByTickers(ctx, tickers)
}

// FindTipRankAssets finds all tip rank assets
func (s *Service) FindTipRankAssets(ctx context.Context) ([]*entities.TipRankAsset, error) {
	s.log.Info(ctx, "find all TipRank assets")
	return s.repo.FindTipRankAssets(ctx)
}
