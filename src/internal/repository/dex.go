package repository

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/subhadip0x539/arc-dex-core-srv/src/internal/domain"
	"github.com/subhadip0x539/arc-dex-core-srv/src/lib"
)

type DexRepository struct {
	pool *pgxpool.Pool
}

func NewDexRepository(pool *pgxpool.Pool) *DexRepository {
	return &DexRepository{pool: pool}
}

func (r *DexRepository) GetDex(ctx context.Context, id string) (*domain.Dex, *domain.Error) {
	d := &domain.Dex{}

	err := r.pool.QueryRow(
		ctx,
		`
		SELECT id, name, color
		FROM dexs
		WHERE id = $1
		`,
		id,
	).Scan(
		&d.ID,
		&d.Name,
		&d.Color,
	)
	if err != nil {
		return nil, domain.NewInternalServerError(err)
	}

	return d, nil
}

func (r *DexRepository) GetDexs(ctx context.Context, query *domain.Query) ([]*domain.Dex, *domain.Error) {
	condition, args, err := domain.NewCondition("d", query.Sort, query.Filters)
	if err != nil {
		return nil, domain.NewInternalServerError(err)
	}

	rows, err := r.pool.Query(ctx, fmt.Sprintf(`
    	SELECT d.id, d.name, d.color
    	FROM dexs d
    	%s
	`, condition), args...)
	if err != nil {
		return nil, domain.NewInternalServerError(err)
	}
	defer rows.Close()

	dexs := []*domain.Dex{}

	for rows.Next() {
		d := &domain.Dex{}

		if err := rows.Scan(
			&d.ID,
			&d.Name,
			&d.Color,
		); err != nil {
			return nil, domain.NewInternalServerError(err)
		}

		dexs = append(dexs, d)
	}

	return dexs, nil
}

func (r *DexRepository) CreateDex(ctx context.Context, dex domain.Dex) (*domain.Dex, *domain.Error) {
	id := lib.Slugify(dex.Name)

	row := r.pool.QueryRow(ctx, `
		INSERT INTO dexs (id, name, color)
		VALUES ($1, $2, $3)
		RETURNING id, name, color
	`,
		id,
		dex.Name,
		dex.Color,
	)

	var d domain.Dex

	err := row.Scan(
		&d.ID,
		&d.Name,
		&d.Color,
	)
	if err != nil {
		return nil, domain.NewInternalServerError(err)
	}

	return &d, nil
}
