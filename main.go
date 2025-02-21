package main

import (
	"bytes"
	"context"
	"io/ioutil"
	"log"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/s3"
)

const bucket = "testbucketcorrectme"
const key = "example.txt"

func main() {
	data, err := ioutil.ReadFile("example.txt")
	if err != nil {
		log.Fatalf("error while reading file: %v", err)
	}

	cfg, err := config.LoadDefaultConfig(context.TODO(), config.WithRegion("ru-central1"))
	if err != nil {
		log.Fatalf("error while configuring S3: %v", err)
	}

	s3Client := s3.NewFromConfig(cfg)

	result, err := s3Client.PutObject(context.TODO(), &s3.PutObjectInput{
		Bucket:      aws.String(bucket),
		Key:         aws.String(key),
		Body:        bytes.NewReader(data),
		ContentType: aws.String("text/plain"),
	})

	if err != nil {
		log.Fatalf("error while sending object to bucket: %v", err)
	}

	log.Println("object sent successfully, size:", result.Size)
}
