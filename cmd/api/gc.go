package main

import (
	"context"
	"log"
	"minio-go-s3/internal/repository"
	"minio-go-s3/internal/storage"
	"time"
)

type GarbageCollector struct {
	repo     repository.MetaRepository
	storage  storage.ObjectStore
	interval time.Duration
}

func NewGarbageCollector(repo repository.MetaRepository, storage storage.ObjectStore, interval time.Duration) *GarbageCollector {
	return &GarbageCollector{
		repo:     repo,
		storage:  storage,
		interval: interval,
	}
}

func (gc *GarbageCollector) Start(ctx context.Context) {
	ticker := time.NewTicker(gc.interval)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			log.Println("GC stopped")
			return
		case <-ticker.C:
			gc.runCleanup(ctx)
		}
	}
}

func (gc *GarbageCollector) runCleanup(ctx context.Context) {
	log.Println("Starting expired paste cleanup")

	expired, err := gc.repo.ListExpired(ctx, time.Now())
	if err != nil {
		log.Printf("Failed to list expired pastes: %v", err)
		return
	}

	for _, paste := range expired {
		// Delete from storage first
		if err := gc.storage.Delete(ctx, paste.ID); err != nil {
			log.Printf("Failed to delete paste %s from storage: %v", paste.ID, err)
			continue
		}

		// Then delete from repository
		if err := gc.repo.Delete(ctx, paste.ID); err != nil {
			log.Printf("Failed to delete paste %s from repository: %v", paste.ID, err)
			// Note: Storage is already deleted; manual recovery might be needed
		} else {
			log.Printf("Deleted expired paste: %s", paste.ID)
		}
	}

	log.Println("Expired paste cleanup completed")
}
