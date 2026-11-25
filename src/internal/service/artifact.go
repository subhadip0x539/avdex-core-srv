package service

import (
	"context"
	"fmt"

	"github.com/subhadip0x539/arc-dex-core-srv/src/internal/domain"
	"github.com/subhadip0x539/arc-dex-core-srv/src/internal/dto"
	"github.com/subhadip0x539/arc-dex-core-srv/src/internal/port"
)

type ArtifactService struct {
	repo port.ArtifactRepository
}

func NewArtifactService(repository port.ArtifactRepository) *ArtifactService {
	return &ArtifactService{repo: repository}
}

func (s *ArtifactService) GetArtifact(ctx context.Context, id string) (*domain.Artifact, *domain.Error) {
	result, err := s.repo.GetArtifact(ctx, id)
	if err != nil {
		return nil, err
	}

	return result, nil
}

func (s *ArtifactService) GetArtifacts(ctx context.Context, query dto.GetArtifactsRequestBody) ([]*domain.Artifact, *domain.Error) {
	qos := domain.NewQuery(query.Sort, query.Filters)

	result, err := s.repo.GetArtifacts(ctx, qos)
	if err != nil {
		return nil, err
	}

	return result, nil
}

func (s *ArtifactService) CreateArtifact(ctx context.Context, payload dto.CreateArtifactRequestBody) (*domain.Artifact, *domain.Error) {
	artifact := domain.Artifact{
		Name:      payload.Name,
		Image:     payload.Image,
		Reference: payload.Reference,
		Dex:       domain.Dex{ID: payload.DexID},
		Extra:     payload.Extra,
	}

	fmt.Print(artifact)

	result, err := s.repo.CreateArtifact(ctx, artifact)
	if err != nil {
		return nil, err
	}

	sites := make([]domain.Site, len(payload.Sites))

	for i, site := range payload.Sites {
		sites[i] = domain.Site{
			Name: site.Name,
			URL:  site.URL,
		}
	}

	rs, err := s.repo.CreateArtifactSites(ctx, result.ID, sites)
	if err != nil {
		return nil, err
	}
	result.Sites = rs

	return result, nil
}
