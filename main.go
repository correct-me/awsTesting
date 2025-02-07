package main

import (
	"bytes"
	"context"
	"io"
	"log"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/s3"
)

func main() {

	bucket := "testbucketcorrectme"
	key := "IMG_5300.MP4"

	cfg, err := config.LoadDefaultConfig(context.TODO(), config.WithRegion("us-west-2"))
	if err != nil {
		log.Fatalf("error while conf s3: %v", err)
	}

	s3Client := s3.NewFromConfig(cfg)

	result, err := s3Client.GetObject(context.TODO(), &s3.GetObjectInput{
		Bucket: aws.String(bucket),
		Key:    aws.String(key),
	})

	defer result.Body.Close()

	buf := bytes.NewBuffer(nil)

	if _, err := io.Copy(buf, result.Body); err != nil {
		log.Fatal(err)
	}

	log.Println(string(buf.Bytes()))

}
