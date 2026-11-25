package port

import (
	"context"

	"github.com/subhadip0x539/arc-dex-core-srv/src/internal/domain"
	"github.com/subhadip0x539/arc-dex-core-srv/src/internal/dto"
)

type DexRepository interface {
	GetDex(ctx context.Context, id string) (*domain.Dex, *domain.Error)
	GetDexs(ctx context.Context, query *domain.Query) ([]*domain.Dex, *domain.Error)
	CreateDex(ctx context.Context, dex domain.Dex) (*domain.Dex, *domain.Error)
}

type DexService interface {
	GetDex(ctx context.Context, id string) (*domain.Dex, *domain.Error)
	GetDexs(ctx context.Context, query dto.GetDexsRequestBody) ([]*domain.Dex, *domain.Error)
	CreateDex(ctx context.Context, req dto.CreateDexRequestBody) (*domain.Dex, *domain.Error)
}
