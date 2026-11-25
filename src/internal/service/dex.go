package service

import (
	"context"

	"github.com/subhadip0x539/arc-dex-core-srv/src/internal/domain"
	"github.com/subhadip0x539/arc-dex-core-srv/src/internal/dto"
	"github.com/subhadip0x539/arc-dex-core-srv/src/internal/port"
)

type DexService struct {
	repo port.DexRepository
}

func NewDexService(repository port.DexRepository) *DexService {
	return &DexService{repo: repository}
}

func (s *DexService) GetDex(ctx context.Context, id string) (*domain.Dex, *domain.Error) {
	result, err := s.repo.GetDex(ctx, id)
	if err != nil {
		return nil, err
	}

	return result, nil
}

func (s *DexService) GetDexs(ctx context.Context, query dto.GetDexsRequestBody) ([]*domain.Dex, *domain.Error) {
	qos := domain.NewQuery(query.Sort, query.Filters)

	result, err := s.repo.GetDexs(ctx, qos)
	if err != nil {
		return nil, err
	}

	return result, nil
}

func (s *DexService) CreateDex(ctx context.Context, payload dto.CreateDexRequestBody) (*domain.Dex, *domain.Error) {
	dex := domain.Dex{
		Name:  payload.Name,
		Color: payload.Color,
	}

	result, err := s.repo.CreateDex(ctx, dex)
	if err != nil {
		return nil, err
	}

	return result, nil
}
