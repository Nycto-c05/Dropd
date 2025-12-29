package repository

import (
	"context"
	"time"
)

type PasteMeta struct {
	ID string `json:"id"`
	// UUID           string    `json:"uuid"`
	IdempotencyKey string    `json:"idempotency_key"`
	Filename       string    `json:"filename"`
	CreatedAt      time.Time `json:"created_at"`
	ExpiresAt      time.Time `json:"expires_at"`
}

type MetaRepository interface {
	GetByIdempotencyKey(ctx context.Context, key string) (*PasteMeta, error)
	GetByID(ctx context.Context, id string) (*PasteMeta, error)
	Insert(ctx context.Context, meta *PasteMeta) (*PasteMeta, error)
	ListExpired(ctx context.Context, now time.Time) ([]*PasteMeta, error)
	Delete(ctx context.Context, id string) error
}
