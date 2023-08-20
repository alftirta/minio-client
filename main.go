package main

import (
	"context"
	"log"
	"os"
	"time"

	"github.com/joho/godotenv"
	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
)

func init() {
	// load .env in the current path
	if err := godotenv.Load(); err != nil {
		log.Fatalf(".env not found, err: %q", err)
	}
}

func getEnvValue(key string) string {
	value, ok := os.LookupEnv(key)
	if !ok {
		log.Fatalf("missing %s in .env", key)
	}
	return value
}

func main() {
	endpoint := getEnvValue("MINIO_ENDPOINT")
	accessKey := getEnvValue("MINIO_ACCESS_KEY")
	secretKey := getEnvValue("MINIO_SECRET_KEY")
	useSSL := false
	if getEnvValue("ENVIRONMENT") == "production" {
		useSSL = true
	}

	// initialize minio client
	minioClient, err := minio.New(endpoint, &minio.Options{
		Creds:  credentials.NewStaticV4(accessKey, secretKey, ""),
		Secure: useSSL,
	})
	if err != nil {
		log.Fatalf("minio client instantiation failed, err: %q", err)
	}

	// get total buckets
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	buckets, err := minioClient.ListBuckets(ctx)
	if err != nil {
		log.Fatalf("list all buckets failed, err: %q", err)
	}
	log.Printf("Total Buckets: %d", len(buckets))
}
