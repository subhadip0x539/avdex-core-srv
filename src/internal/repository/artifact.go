package repository

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/subhadip0x539/arc-dex-core-srv/src/internal/domain"
)

type ArtifactRepository struct {
	pool *pgxpool.Pool
}

func NewArtifactRepository(pool *pgxpool.Pool) *ArtifactRepository {
	return &ArtifactRepository{pool: pool}
}

func (r *ArtifactRepository) GetArtifact(ctx context.Context, id string) (*domain.Artifact, *domain.Error) {
	a, e := &domain.Artifact{}, []byte{}

	err := r.pool.QueryRow(
		ctx,
		`
		SELECT id, name, image, reference, dex_id, extra
		FROM artifacts
		WHERE id = $1
		`,
		id,
	).Scan(
		&a.ID,
		&a.Name,
		&a.Image,
		&a.Reference,
		&a.Dex.ID,
		&e,
	)
	if err != nil {
		return nil, domain.NewInternalServerError(err)
	}

	err = r.pool.QueryRow(ctx, `
		SELECT id, name, color
		FROM dexs
		WHERE id = $1
	`, a.Dex.ID).Scan(
		&a.Dex.ID,
		&a.Dex.Name,
		&a.Dex.Color,
	)
	if err != nil {
		return nil, domain.NewInternalServerError(err)
	}

	rows, err := r.pool.Query(ctx, `
        SELECT id, name, url
        FROM artifact_sites
        WHERE artifact_id = $1
    `, id)
	if err != nil {
		return nil, domain.NewInternalServerError(err)
	}
	defer rows.Close()

	for rows.Next() {
		var s domain.Site
		if err := rows.Scan(&s.ID, &s.Name, &s.URL); err != nil {
			return nil, domain.NewInternalServerError(err)
		}

		a.Sites = append(a.Sites, s)
	}
	if err := rows.Err(); err != nil {
		return nil, domain.NewInternalServerError(err)
	}

	if err := json.Unmarshal(e, &a.Extra); err != nil {
		return nil, domain.NewInternalServerError(err)
	}

	return a, nil
}

func (r *ArtifactRepository) GetArtifacts(ctx context.Context, query *domain.Query) ([]*domain.Artifact, *domain.Error) {
	condition, args, err := domain.NewCondition("a", query.Sort, query.Filters)
	if err != nil {
		return nil, domain.NewInternalServerError(err)
	}

	rows, err := r.pool.Query(ctx, fmt.Sprintf(`
		SELECT a.id, a.name, COALESCE(a.image, ''), a.reference, d.id, d.name, d.color, a.extra
		FROM artifacts a
		LEFT JOIN dexs d ON d.id = a.dex_id
		%s
	`, condition), args...)
	if err != nil {
		return nil, domain.NewInternalServerError(err)
	}
	defer rows.Close()

	amap, ids := map[string]*domain.Artifact{}, []string{}

	for rows.Next() {
		a, e := &domain.Artifact{}, []byte{}

		if err := rows.Scan(
			&a.ID,
			&a.Name,
			&a.Image,
			&a.Reference,
			&a.Dex.ID,
			&a.Dex.Name,
			&a.Dex.Color,
			&e,
		); err != nil {
			return nil, domain.NewInternalServerError(err)
		}

		if err := json.Unmarshal(e, &a.Extra); err != nil {
			return nil, domain.NewInternalServerError(err)
		}

		amap[a.ID] = a
		ids = append(ids, a.ID)
	}

	srows, err := r.pool.Query(ctx, `
        SELECT id, name, artifact_id, url
        FROM artifact_sites
        WHERE artifact_id = ANY($1)
    `, ids)
	if err != nil {
		return nil, domain.NewInternalServerError(err)
	}
	defer srows.Close()

	for srows.Next() {
		var id string
		var s domain.Site

		if err := srows.Scan(&id, &s.ID, &s.Name, &s.URL); err != nil {
			return nil, domain.NewInternalServerError(err)
		}

		if a, ok := amap[id]; ok {
			a.Sites = append(a.Sites, s)
		}
	}

	if err := srows.Err(); err != nil {
		return nil, domain.NewInternalServerError(err)
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
		return nil, domain.NewInternalServerError(err)
	}

	id := uuid.NewString()

	row := r.pool.QueryRow(ctx, `
		INSERT INTO artifacts (id, name, image, reference, dex_id, extra)
		VALUES ($1, $2, $3, $4, $5, $6)
		RETURNING id, name, image, reference, dex_id, extra
	`,
		id,
		artifact.Name,
		artifact.Image,
		artifact.Reference,
		artifact.Dex.ID,
		extra,
	)

	var a domain.Artifact
	var e []byte

	err = row.Scan(
		&a.ID,
		&a.Name,
		&a.Image,
		&a.Reference,
		&a.Dex.ID,
		&e,
	)
	if err != nil {
		return nil, domain.NewInternalServerError(err)
	}

	if e != nil {
		if err := json.Unmarshal(e, &a.Extra); err != nil {
			return nil, domain.NewInternalServerError(err)
		}
	}

	return &a, nil
}

func (r *ArtifactRepository) CreateArtifactSites(ctx context.Context, artifactID string, sites []domain.Site) ([]domain.Site, *domain.Error) {
	created := make([]domain.Site, 0, len(sites))

	for _, site := range sites {
		id := uuid.NewString()

		var s domain.Site

		err := r.pool.QueryRow(ctx, `
            INSERT INTO artifact_sites (
                id,
                name,
                artifact_id,
                url
            )
            VALUES ($1, $2, $3, $4)
            RETURNING id, name, url
        `,
			id,
			site.Name,
			artifactID,
			site.URL,
		).Scan(
			&s.ID,
			&s.Name,
			&s.URL,
		)
		if err != nil {
			return nil, domain.NewInternalServerError(err)
		}

		created = append(created, s)
	}

	return created, nil
}
