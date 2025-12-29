package main

import (
	"context"
	"fmt"
	"io"
	"log"
	"os"

	minio "github.com/minio/minio-go/v7"
	credentials "github.com/minio/minio-go/v7/pkg/credentials"
)

var print = fmt.Println

func CreateBucket(ctx context.Context, client *minio.Client, bucketName, location string) error {
	err := client.MakeBucket(ctx, bucketName, minio.MakeBucketOptions{
		Region: location,
	})

	if err == nil {
		log.Printf("Successfully created %s\n", bucketName)
		return nil
	}

	// check if bucket already exists
	exists, errBucketExists := client.BucketExists(ctx, bucketName)
	if errBucketExists != nil {
		return fmt.Errorf("check bucket exists failed: %w", errBucketExists)
	}

	if exists {
		log.Printf("Bucket %s already exists\n", bucketName)
		return nil // important: not an error
	}

	return fmt.Errorf("create bucket %s failed: %w", bucketName, err)
}

func PutObject(ctx context.Context, client *minio.Client, filepath string, bucketName string, objName string) error {
	//os.File is an impl of interface io.Reader interface
	file, err := os.Open(filepath)
	if err != nil {
		print(err)
		return err
	}
	defer file.Close()

	fileStat, err := file.Stat()
	if err != nil {
		print(err)
		return err
	}

	uploadInfo, err := client.PutObject(ctx, bucketName, objName, file, fileStat.Size(), minio.PutObjectOptions{ContentType: "application/octet-stream"})
	if err != nil {
		print(err)
		return err
	}

	print("Succesfully uploaded bytes: ", uploadInfo.Size)
	return nil
}

func GetObject(ctx context.Context, client *minio.Client, bucketName string, objectName string) ([]byte, error) {
	object, err := client.GetObject(ctx, bucketName, objectName, minio.GetObjectOptions{})
	if err != nil {
		log.Println("error while fetching the object "+objectName+" from the bucket "+bucketName+" :", err)
		return nil, err
	}
	defer object.Close()

	data, err := io.ReadAll(object)
	if err != nil {
		log.Println("error reading data from object")
		return nil, err
	}
	return data, nil
}

func main() {
	endpoint := "localhost:9000"
	accessKeyID := "nycto" //TODO: Get from Env
	secretAccessKey := "nycto1234"
	useSSL := false

	// Initialize minio client object.
	minioClient, err := minio.New(endpoint, &minio.Options{
		Creds:  credentials.NewStaticV4(accessKeyID, secretAccessKey, ""),
		Secure: useSSL,
	})
	if err != nil {
		log.Fatalln(err)
	}

	log.Printf("%v\n", minioClient) // minioClient is now setup

	bucketName := "pastetin"
	location := "us-east-1"

	if err := CreateBucket(context.Background(), minioClient, bucketName, location); err != nil {
		log.Println("error creating bucket: ", err)
	}

	if err := PutObject(context.Background(), minioClient, "ipsum.txt", bucketName, "go_mod.txt"); err != nil {
		log.Println("error uploading object: ", err)
	}

	data, err := GetObject(context.Background(), minioClient, bucketName, "go_mod.txt")
	if err != nil {
		log.Println("error reading bytes from object: ", err)
	}
	print(string(data))

}

