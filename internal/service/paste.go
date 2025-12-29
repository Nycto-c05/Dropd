package service

import (
	"context"
	"errors"
	"time"

	"minio-go-s3/internal/idgen"
	"minio-go-s3/internal/repository"
	"minio-go-s3/internal/storage"
)

var (
	ErrExpiredPaste = errors.New("paste expired")
	ErrEmptyPaste   = errors.New("empty paste")
)

type PasteService interface {
	CreatePaste(
		ctx context.Context,
		content []byte,
		filename string,
		idempotencyKey string,
		ttl time.Duration,
	) (*repository.PasteMeta, error)

	GetPaste(
		ctx context.Context,
		id string,
	) (*repository.PasteMeta, []byte, error)

	DeletePaste(ctx context.Context, key string) error
}

type pasteService struct {
	metaRepo repository.MetaRepository
	store    storage.ObjectStore
	idgen    idgen.IDGenerator
}

func NewPasteService(
	metaRepo repository.MetaRepository,
	store storage.ObjectStore,
	idgen idgen.IDGenerator,
) PasteService {
	return &pasteService{
		metaRepo: metaRepo,
		store:    store,
		idgen:    idgen,
	}
}

func (s *pasteService) CreatePaste(
	ctx context.Context,
	content []byte,
	filename string,
	idempotencyKey string,
	ttl time.Duration,
) (*repository.PasteMeta, error) {

	//0. if empty paste
	if len(content) == 0 {
		return nil, ErrEmptyPaste
	}

	// 1. Idempotency check
	if idempotencyKey != "" {
		existing, err := s.metaRepo.GetByIdempotencyKey(ctx, idempotencyKey)
		if err != nil {
			return nil, err
		}
		if existing != nil {
			return existing, nil
		}
	}

	// 2. Generate short ID (Snowflake-style → base62)
	rawID := s.idgen.Next()
	id := idgen.EncodeBase62(rawID)

	// 3. Write object first (S3 / MinIO)
	if err := s.store.Put(ctx, id, content); err != nil {
		return nil, err
	}

	now := time.Now()

	meta := &repository.PasteMeta{
		ID:             id,
		IdempotencyKey: idempotencyKey,
		Filename:       filename,
		CreatedAt:      now,
		ExpiresAt:      now.Add(ttl),
	}

	// 4. Persist metadata
	_, err := s.metaRepo.Insert(ctx, meta)
	if err != nil {
		// compensate
		_ = s.store.Delete(ctx, id)
		return nil, err
	}

	return meta, nil
}

func (s *pasteService) GetPaste(
	ctx context.Context,
	id string,
) (*repository.PasteMeta, []byte, error) {

	meta, err := s.metaRepo.GetByID(ctx, id)
	if err != nil {
		return nil, nil, err
	}
	if meta == nil {
		return nil, nil, nil
	}

	if meta.ExpiresAt.Before(time.Now()) {
		return nil, nil, ErrExpiredPaste
	}

	data, err := s.store.Get(ctx, id)
	if err != nil {
		return nil, nil, err
	}

	return meta, data, nil
}

func (s *pasteService) DeletePaste(ctx context.Context, id string) error {
	if err := s.store.Delete(ctx, id); err != nil {
		return err
	}

	if err := s.metaRepo.Delete(ctx, id); err != nil {
		return err
	}

	return nil
}
