package repository

import (
	"context"
	"encoding/json"
	"fmt"
	"log/slog"
	"net/http"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/subhadip0x539/arc-dex-core-srv/src/internal/domain"
)

const defaultSiteID = "JavTrailers"

type ArtifactRepository struct {
	pool *pgxpool.Pool
}

func NewArtifactRepository(pool *pgxpool.Pool) *ArtifactRepository {
	return &ArtifactRepository{pool: pool}
}

func (r *ArtifactRepository) GetArtifact(ctx context.Context, id string) (*domain.Artifact, *domain.Error) {
	a := &domain.Artifact{}
	var e []byte

	err := r.pool.QueryRow(
		ctx,
		`
		SELECT id, name, gender, height, birthdate, image, chest, hips, waist, element_id, extra
		FROM artifacts
		WHERE id = $1
		`,
		id,
	).Scan(
		&a.ID,
		&a.Name,
		&a.Gender,
		&a.Height,
		&a.Birthdate,
		&a.Image,
		&a.Measurements.Chest,
		&a.Measurements.Hips,
		&a.Measurements.Waist,
		&a.Element,
		&e,
	)
	if err != nil {
		slog.Error(err.Error())
		return nil, &domain.Error{Code: http.StatusInternalServerError, Message: string(domain.ErrorMessageInternalServerError)}
	}

	if err := json.Unmarshal(e, &a.Extra); err != nil {
		slog.Error(err.Error())
		return nil, &domain.Error{Code: http.StatusInternalServerError, Message: string(domain.ErrorMessageInternalServerError), Trace: err}
	}

	rows, err := r.pool.Query(
		ctx,
		`
		SELECT a_s.site_id, s.name, a_s.url
		FROM artifact_sites AS a_s
		JOIN sites AS s ON s.id = a_s.site_id
		WHERE a_S.artifact_id = $1
		`,
		id,
	)
	if err != nil {
		slog.Error(err.Error())
		return nil, &domain.Error{Code: http.StatusInternalServerError, Message: string(domain.ErrorMessageInternalServerError), Trace: err}
	}
	defer rows.Close()

	for rows.Next() {
		var s domain.Site
		if err := rows.Scan(&s.ID, &s.Name, &s.URL); err != nil {
			slog.Error(err.Error())
			return nil, &domain.Error{Code: http.StatusInternalServerError, Message: string(domain.ErrorMessageInternalServerError), Trace: err}
		}

		a.Sites = append(a.Sites, s)

		if s.ID == defaultSiteID {
			sc := s
			a.Domain = &sc
		}
	}
	if err := rows.Err(); err != nil {
		slog.Error(err.Error())
		return nil, &domain.Error{Code: http.StatusInternalServerError, Message: string(domain.ErrorMessageInternalServerError), Trace: err}
	}

	now := time.Now()
	age := now.Year() - a.Birthdate.Year()

	a.Age = age

	return a, nil
}

func (r *ArtifactRepository) GetArtifacts(ctx context.Context, query domain.Query) ([]*domain.Artifact, *domain.Error) {
	rows, err := r.pool.Query(ctx, fmt.Sprintf(`
		SELECT id, name, gender, height, birthdate, image, chest, hips, waist, element_id, extra
		FROM artifacts
		ORDER BY %s %s
	`, query.Sort.By, query.Sort.Order))
	if err != nil {
		slog.Error(err.Error())
		return nil, &domain.Error{Code: http.StatusInternalServerError, Message: string(domain.ErrorMessageInternalServerError), Trace: err}
	}
	defer rows.Close()

	amap, ids := map[string]*domain.Artifact{}, []string{}

	for rows.Next() {
		a, e := &domain.Artifact{}, []byte{}

		if err := rows.Scan(
			&a.ID,
			&a.Name,
			&a.Gender,
			&a.Height,
			&a.Birthdate,
			&a.Image,
			&a.Measurements.Chest,
			&a.Measurements.Hips,
			&a.Measurements.Waist,
			&a.Element,
			&e,
		); err != nil {
			slog.Error(err.Error())
			return nil, &domain.Error{Code: http.StatusInternalServerError, Message: string(domain.ErrorMessageInternalServerError), Trace: err}
		}

		if err := json.Unmarshal(e, &a.Extra); err != nil {
			slog.Error(err.Error())
			return nil, &domain.Error{Code: http.StatusInternalServerError, Message: string(domain.ErrorMessageInternalServerError), Trace: err}
		}

		now := time.Now()
		age := now.Year() - a.Birthdate.Year()

		if now.Month() < a.Birthdate.Month() ||
			(now.Month() == a.Birthdate.Month() && now.Day() < a.Birthdate.Day()) {
			age--
		}

		a.Age = age

		amap[a.ID] = a
		ids = append(ids, a.ID)
	}

	srows, err := r.pool.Query(ctx, `
		SELECT a_s.artifact_id, a_s.site_id, s.name, a_s.url
		FROM artifact_sites a_s
		JOIN sites s ON s.id = a_s.site_id
		WHERE a_s.artifact_id = ANY($1)
	`, ids)
	if err != nil {
		slog.Error(err.Error())
		return nil, &domain.Error{Code: http.StatusInternalServerError, Message: string(domain.ErrorMessageInternalServerError), Trace: err}
	}
	defer srows.Close()

	for srows.Next() {
		var id string
		var s domain.Site

		if err := srows.Scan(&id, &s.ID, &s.Name, &s.URL); err != nil {
			slog.Error(err.Error())
			return nil, &domain.Error{Code: http.StatusInternalServerError, Message: string(domain.ErrorMessageInternalServerError), Trace: err}
		}

		fmt.Println(s.ID)

		if a, ok := amap[id]; ok {
			a.Sites = append(a.Sites, s)

			if s.ID == defaultSiteID {
				sc := s
				a.Domain = &sc
			}
		}
	}

	if err := srows.Err(); err != nil {
		slog.Error(err.Error())
		return nil, &domain.Error{Code: http.StatusInternalServerError, Message: string(domain.ErrorMessageInternalServerError), Trace: err}
	}

	artifacts := make([]*domain.Artifact, 0, len(ids))
	for _, id := range ids {
		artifacts = append(artifacts, amap[id])
	}

	return artifacts, nil
}

func (r *ArtifactRepository) CreateArtifact(ctx context.Context, artifact domain.Artifact) (*domain.Artifact, *domain.Error) {
	extra, err := json.Marshal(artifact.Extra)
	if err != nil {
		slog.Error(err.Error())
		return nil, &domain.Error{Code: http.StatusInternalServerError, Message: string(domain.ErrorMessageInternalServerError), Trace: err}
	}

	id := uuid.New()

	row := r.pool.QueryRow(ctx, `
		INSERT INTO artifacts (id, name, gender, height, birthdate, image, chest, hips, waist, element_id, extra)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11)
		RETURNING id, name, gender, height, birthdate, image, chest, hips, waist, element_id, extra
	`,
		id.String(),
		artifact.Name,
		artifact.Gender,
		artifact.Height,
		artifact.Birthdate,
		artifact.Image,
		artifact.Measurements.Chest,
		artifact.Measurements.Hips,
		artifact.Measurements.Waist,
		artifact.Element,
		extra,
	)

	var a domain.Artifact
	var e []byte

	err = row.Scan(
		&a.ID,
		&a.Name,
		&a.Gender,
		&a.Height,
		&a.Birthdate,
		&a.Image,
		&a.Measurements.Chest,
		&a.Measurements.Hips,
		&a.Measurements.Waist,
		&a.Element,
		&e,
	)
	if err != nil {
		slog.Error(err.Error())
		return nil, &domain.Error{Code: http.StatusInternalServerError, Message: string(domain.ErrorMessageInternalServerError), Trace: err}
	}

	if e != nil {
		if err := json.Unmarshal(e, &a.Extra); err != nil {
			slog.Error(err.Error())
			return nil, &domain.Error{Code: http.StatusInternalServerError, Message: string(domain.ErrorMessageInternalServerError), Trace: err}
		}
	}

	return &a, nil
}

func (r *ArtifactRepository) CreateArtifactSites(ctx context.Context, artifactID string, sites []domain.Site) ([]domain.Site, *domain.Error) {
	ptx, err := r.pool.Begin(ctx)
	if err != nil {
		slog.Error(err.Error())
		return nil, &domain.Error{Code: http.StatusInternalServerError, Message: string(domain.ErrorMessageInternalServerError), Trace: err}
	}
	defer ptx.Rollback(ctx)

	var result []domain.Site

	for _, s := range sites {
		row := ptx.QueryRow(ctx, `
			WITH upserted AS (
				INSERT INTO artifact_sites (artifact_id, site_id, url)
				VALUES ($1, $2, $3)
				ON CONFLICT (artifact_id, site_id) DO UPDATE SET url = EXCLUDED.url
				RETURNING artifact_id, site_id, url
			)
			SELECT s.id, s.name, upserted.url
			FROM upserted
			JOIN sites s ON s.id = upserted.site_id
		`, artifactID, s.ID, s.URL)

		var s domain.Site

		if err := row.Scan(&s.ID, &s.Name, &s.URL); err != nil {
			slog.Error(err.Error())
			return nil, &domain.Error{Code: http.StatusInternalServerError, Message: string(domain.ErrorMessageInternalServerError), Trace: err}
		}

		result = append(result, s)
	}

	if err := ptx.Commit(ctx); err != nil {
		slog.Error(err.Error())
		return nil, &domain.Error{Code: 500, Message: string(domain.ErrorMessageInternalServerError), Trace: err}
	}

	return result, nil
}
