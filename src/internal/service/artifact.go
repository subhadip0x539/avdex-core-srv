package service

import (
	"context"
	"log/slog"
	"net/http"
	"time"

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

func (s *ArtifactService) GetArtifacts(ctx context.Context, query dto.GetArtifactsRequestQuery) ([]*domain.Artifact, *domain.Error) {
	qos := domain.NewQuery(query.Sort)

	result, err := s.repo.GetArtifacts(ctx, qos)
	if err != nil {
		return nil, err
	}

	return result, nil
}

func (s *ArtifactService) CreateArtifact(ctx context.Context, payload dto.CreateArtifactRequestBody) (*domain.Artifact, *domain.Error) {
	var birthdate time.Time

	if payload.Birthdate != "" {
		var err error
		birthdate, err = time.Parse("2006-01-02", payload.Birthdate)
		if err != nil {
			slog.Error(err.Error())
			return nil, &domain.Error{Code: http.StatusBadRequest, Message: string(domain.ErrorMessageBadRequest), Trace: err}
		}
	}

	artifact := domain.Artifact{
		Name:         payload.Name,
		Gender:       domain.Gender(payload.Gener),
		Height:       payload.Height,
		Birthdate:    birthdate,
		Measurements: domain.Measurements(payload.Measurements),
		Image:        payload.Image,
		Element:      payload.Element,
		Extra:        payload.Extra,
	}

	result, err := s.repo.CreateArtifact(ctx, artifact)
	if err != nil {
		return nil, err
	}

	sites := make([]domain.Site, len(payload.Sites))

	for i, site := range payload.Sites {
		sites[i] = domain.Site{
			ID:   site.ID,
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
