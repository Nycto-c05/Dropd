package main

import (
	"context"
	"log"
	"minio-go-s3/internal/blob"
	"minio-go-s3/internal/db"
	"minio-go-s3/internal/idgen"
	"minio-go-s3/internal/repository"
	"minio-go-s3/internal/service"
	"minio-go-s3/internal/storage"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func main() {

	config := &config{
		addr:    ":8083",
		baseURL: "http://localhost:8083",
		db: dbConfig{
			addr:         "postgresql://postgres:postgres@localhost:5432/paste?sslmode=disable",
			maxOpenConns: 30,
			maxIdleConns: 30,
			maxIdleTime:  "15m",
		},
		blob: blobConfig{
			endpoint:        "localhost:9000",
			accessKeyID:     "nycto",
			secretAccessKey: "nycto1234",
			bucket:          "pastes",
			useSSL:          false,
		},
	}

	//connect to db, verify conn, and get client
	dbClient, err := db.NewDb(config.db.addr, config.db.maxOpenConns, config.db.maxIdleConns, config.db.maxIdleTime)
	if err != nil {
		log.Panic(err)
	}
	defer dbClient.Close()
	log.Print("DB connection pool established")

	//connect to s3
	s3Client, err := blob.NewObj(config.blob.endpoint, config.blob.accessKeyID, config.blob.secretAccessKey, config.blob.useSSL)
	if err != nil {
		log.Panic(err)
	}
	log.Print("Object Storage Connnection established")

	//Initialize impl of repo and object stores, and the service
	repo := repository.NewPostgresMetaRepository(dbClient)
	objectStore := storage.NewMinioStore(s3Client, config.blob.bucket)
	idGen := idgen.NewIDGenerator()
	gc := NewGarbageCollector(repo, objectStore, 24*time.Hour)

	pasteSvc := service.NewPasteService(repo, objectStore, idGen)

	//Create application
	app := &application{
		config:   *config,
		pasteSvc: pasteSvc,
		gc:       gc,
	}

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	go app.gc.Start(ctx)

	mux := app.mount()

	log.Fatal(app.run(mux))

	signalChan := make(chan os.Signal, 1)
	signal.Notify(signalChan, syscall.SIGINT, syscall.SIGTERM)
	<-signalChan
	cancel()
	log.Println("Shutting down...")

}
