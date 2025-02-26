package main

import (
	"bytes"
	"context"
	"io"
	"log"
	"os"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/s3"
)

const bucket = "testbucketcorrectme"
const key = "2025-02-06 17-28-35.mp4"

func main() {
	cfg, err := config.LoadDefaultConfig(context.TODO(), config.WithRegion("ru-central1"))
	if err != nil {
		log.Fatalf("error while configuring S3: %v", err)
	}

	s3Client := s3.NewFromConfig(cfg)

	if err := uploadFileToBucket(s3Client, "media/2025-02-06 17-28-35.mp4"); err != nil {
		log.Fatal(err)
	}

	if err := downloadFileFromBucket(s3Client, "media/2025-02-06 17-28-35_COPY.mp4"); err != nil {
		log.Fatal(err)
	}

}

func uploadFileToBucket(client *s3.Client, filepath string) error {
	data, err := os.ReadFile(filepath)
	if err != nil {
		return err
	}

	if _, err := client.PutObject(context.TODO(), &s3.PutObjectInput{
		Bucket:      aws.String(bucket),
		Key:         aws.String(key),
		Body:        bytes.NewReader(data),
		ContentType: aws.String("text/plain"),
	}); err != nil {
		return err
	}

	log.Println("object sent successfully")

	return nil
}

func downloadFileFromBucket(client *s3.Client, dir string) error {
	result, err := client.GetObject(context.TODO(), &s3.GetObjectInput{
		Bucket: aws.String(bucket),
		Key:    aws.String(key),
	})
	if err != nil {
		return err
	}

	objectsBytes, err := io.ReadAll(result.Body)
	if err != nil {
		return err
	}

	file, err := os.Create(dir)
	if err != nil {
		return err
	}
	defer file.Close()

	_, err = file.Write(objectsBytes)
	if err != nil {
		return err
	}

	log.Printf("successfully downloaded data from %s/%s\n to file %s", bucket, key, dir)
	return nil
}
