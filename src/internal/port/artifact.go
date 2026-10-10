package port

import (
	"context"

	"github.com/subhadip0x539/arc-dex-core-srv/src/internal/domain"
	"github.com/subhadip0x539/arc-dex-core-srv/src/internal/dto"
)

type ArtifactRepository interface {
	GetArtifact(ctx context.Context, id string) (*domain.Artifact, *domain.Error)
	GetArtifacts(ctx context.Context, query domain.Query) ([]*domain.Artifact, *domain.Error)
	CreateArtifact(ctx context.Context, artifact domain.Artifact) (*domain.Artifact, *domain.Error)
	CreateArtifactSites(ctx context.Context, artifactID string, sites []domain.Site) ([]domain.Site, *domain.Error)
}

type ArtifactService interface {
	GetArtifact(ctx context.Context, id string) (*domain.Artifact, *domain.Error)
	GetArtifacts(ctx context.Context, query dto.GetArtifactsRequestQuery) ([]*domain.Artifact, *domain.Error)
	CreateArtifact(ctx context.Context, req dto.CreateArtifactRequestBody) (*domain.Artifact, *domain.Error)
}
