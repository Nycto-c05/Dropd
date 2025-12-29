package storage

import (
	"bytes"
	"context"
	"io"

	minio "github.com/minio/minio-go/v7"
)

type MinioStore struct {
	client *minio.Client
	bucket string
}

func NewMinioStore(client *minio.Client, bucket string) *MinioStore {
	return &MinioStore{client: client, bucket: bucket}
}

func (s *MinioStore) Get(ctx context.Context, key string) ([]byte, error) {
	object, err := s.client.GetObject(ctx, s.bucket, key, minio.GetObjectOptions{})
	if err != nil {
		// log.Println("error while fetching the object "+key+" from the bucket "+s.bucket+" :", err)
		return nil, err
	}
	defer object.Close()

	return io.ReadAll(object)
}

func (s *MinioStore) Put(ctx context.Context, key string, data []byte) error {
	reader := bytes.NewReader(data)

	_, err := s.client.PutObject(
		ctx,
		s.bucket,
		key,
		reader,
		int64(len(data)),
		minio.PutObjectOptions{
			ContentType: "text/plain",
		},
	)

	return err
}

func (s *MinioStore) Delete(ctx context.Context, key string) error {
	err := s.client.RemoveObject(ctx, s.bucket, key, minio.RemoveObjectOptions{})

	return err
}
