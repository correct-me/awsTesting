package main

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/s3"
)

const bucket = "testbucketcorrectme"

func main() {
	cfg, err := config.LoadDefaultConfig(context.TODO(), config.WithRegion("ru-central1"))
	if err != nil {
		log.Fatalf("error while configuring S3: %v", err)
	}

	s3Client := s3.NewFromConfig(cfg)

	http.HandleFunc("/upload", uploadVideoHandler(s3Client))
	http.HandleFunc("/download", downloadVideoHandler(s3Client))

	log.Println("Server started at :8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}

func uploadVideoHandler(client *s3.Client) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}

		file, header, err := r.FormFile("file")
		if err != nil {
			http.Error(w, "Error retrieving file from form", http.StatusBadRequest)
			return
		}
		defer file.Close()

		fileBytes, err := io.ReadAll(file)
		if err != nil {
			http.Error(w, "Error reading file", http.StatusInternalServerError)
			return
		}

		key := header.Filename

		_, err = client.PutObject(context.TODO(), &s3.PutObjectInput{
			Bucket:      aws.String(bucket),
			Key:         aws.String(key),
			Body:        bytes.NewReader(fileBytes),
			ContentType: aws.String(header.Header.Get("Content-Type")),
		})
		if err != nil {
			http.Error(w, "Failed to upload file to S3", http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusOK)
		w.Write([]byte("File uploaded successfully"))
	}
}

func downloadVideoHandler(client *s3.Client) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		key := r.URL.Query().Get("key")
		if key == "" {
			http.Error(w, "Missing key parameter", http.StatusBadRequest)
			return
		}

		result, err := client.GetObject(context.TODO(), &s3.GetObjectInput{
			Bucket: aws.String(bucket),
			Key:    aws.String(key),
		})
		if err != nil {
			http.Error(w, "Failed to get file from S3", http.StatusInternalServerError)
			return
		}
		defer result.Body.Close()

		var contentType string
		if result.ContentType != nil {
			contentType = *result.ContentType
		}

		w.Header().Set("Content-Type", contentType)
		w.Header().Set("Content-Disposition", fmt.Sprintf("attachment; filename=%q", key))

		file, err := os.Create("media/new_file.mov")
		if err != nil {
			http.Error(w, "Failed create file", http.StatusInternalServerError)
			return
		}

		if _, err := io.Copy(file, result.Body); err != nil {
			http.Error(w, "Error streaming file", http.StatusInternalServerError)
			return
		}
	}
}
