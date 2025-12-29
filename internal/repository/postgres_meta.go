package repository

import (
	"context"
	"database/sql"
	"time"
)

type PostgresMetaRepository struct {
	db *sql.DB
}

func NewPostgresMetaRepository(db *sql.DB) *PostgresMetaRepository {
	return &PostgresMetaRepository{db: db}
}

func (r *PostgresMetaRepository) GetByIdempotencyKey(ctx context.Context, key string) (*PasteMeta, error) {
	const q = `
		SELECT id, idempotency_key, filename, created_at, expires_at 
		FROM paste_metadata 
		WHERE idempotency_key = $1;
	`
	paste := &PasteMeta{}

	err := r.db.QueryRowContext(ctx, q, key).Scan(
		&paste.ID,
		&paste.IdempotencyKey,
		&paste.Filename,
		&paste.CreatedAt,
		&paste.ExpiresAt,
	)

	if err == sql.ErrNoRows {
		return nil, nil
	}

	if err != nil {
		return nil, err
	}

	return paste, nil
}

func (r *PostgresMetaRepository) GetByID(ctx context.Context, id string) (*PasteMeta, error) {
	const q = `
		SELECT id, idempotency_key, filename, created_at, expires_at 
		FROM paste_metadata 
		WHERE id = $1;
	`
	paste := &PasteMeta{}

	err := r.db.QueryRowContext(ctx, q, id).Scan(
		&paste.ID,
		&paste.IdempotencyKey,
		&paste.Filename,
		&paste.CreatedAt,
		&paste.ExpiresAt,
	)

	if err == sql.ErrNoRows {
		return nil, nil
	}

	if err != nil {
		return nil, err
	}

	return paste, nil

}

func (r *PostgresMetaRepository) Insert(ctx context.Context, meta *PasteMeta) (*PasteMeta, error) {
	const q = `
        INSERT INTO paste_metadata (
			id,
            idempotency_key,
            filename,
            created_at,
            expires_at
        )
        VALUES ($1, $2, $3, $4, $5) RETURNING id;
    `
	inserted := &PasteMeta{}

	err := r.db.QueryRowContext(ctx, q,
		meta.ID,
		meta.IdempotencyKey,
		meta.Filename,
		meta.CreatedAt,
		meta.ExpiresAt,
	).Scan(
		&inserted.ID,
	)

	if err != nil {
		return nil, err
	}

	return inserted, nil
}

func (r *PostgresMetaRepository) ListExpired(
	ctx context.Context,
	now time.Time,
) ([]*PasteMeta, error) {

	const q = `
		SELECT id, idempotency_key, filename, created_at, expires_at
		FROM paste_metadata
		WHERE expires_at <= $1;
	`

	rows, err := r.db.QueryContext(ctx, q, now)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var pastes []*PasteMeta

	for rows.Next() {
		p := &PasteMeta{}
		if err := rows.Scan(
			&p.ID,
			&p.IdempotencyKey,
			&p.Filename,
			&p.CreatedAt,
			&p.ExpiresAt,
		); err != nil {
			return nil, err
		}
		pastes = append(pastes, p)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return pastes, nil
}

func (r *PostgresMetaRepository) Delete(
	ctx context.Context,
	id string,
) error {

	const q = `
		DELETE FROM paste_metadata
		WHERE id = $1;
	`

	res, err := r.db.ExecContext(ctx, q, id)
	if err != nil {
		return err
	}

	rows, err := res.RowsAffected()
	if err != nil {
		return err
	}

	if rows == 0 {
		return sql.ErrNoRows
	}

	return nil
}
