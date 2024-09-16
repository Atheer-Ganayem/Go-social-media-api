package utils

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"mime/multipart"
	"os"
	"path/filepath"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/google/uuid"
)

var S3Client *s3.Client

func InitAWS() {
	cfg, err := config.LoadDefaultConfig(context.TODO(),
		config.WithRegion("eu-central-1"),
		config.WithCredentialsProvider(credentials.NewStaticCredentialsProvider(
			os.Getenv("AWS_ACCESS_KEY_ID"),
			os.Getenv("AWS_SECRET_ACCESS_KEY"),
			"",
		)),)

	if err != nil {
		panic(err)
	}

	S3Client = s3.NewFromConfig(cfg)
}

func UploadFileToS3(file multipart.File, fileHeader *multipart.FileHeader) (string, error) {
	bucketName := os.Getenv("AWS_BUCKET_NAME")
	fileKey := fmt.Sprintf("go-social-media/%s%s", uuid.New().String(), filepath.Ext(fileHeader.Filename))
	fileContent, err := io.ReadAll(file)
	if err != nil {
		return "", err
	}

	_, err = S3Client.PutObject(context.TODO(), &s3.PutObjectInput{
		Bucket: aws.String(bucketName),
		Key: aws.String(fileKey),
		Body: bytes.NewReader(fileContent),
		ContentType: aws.String(fileHeader.Header.Get("Content-Type")),
	})

	if err != nil {
		return "", err
	}

	return fileKey, nil
}

func RemoveFileFromS3(filePath string) error {
	bucketName := os.Getenv("AWS_BUCKET_NAME")

	_, err := S3Client.DeleteObject(context.TODO(), &s3.DeleteObjectInput{
		Bucket: aws.String(bucketName),
		Key: aws.String(filePath),
	})

	return err
}