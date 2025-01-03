package externalServices

import (
	"bytes"
	"context"
	"log"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/s3"
)

type S3Client interface {
	UploadToS3(ctx context.Context, bucketName, key string, content []byte) error
}

type s3Client struct{ Client *s3.Client }

func NewS3Client() (S3Client, error) {

	cfg, err := config.LoadDefaultConfig(context.TODO())
	if err != nil {
		return nil, err
	}
	client := s3.NewFromConfig(cfg)
	return &s3Client{Client: client}, nil
}

func (s *s3Client) UploadToS3(ctx context.Context, bucketName, key string, content []byte) error {
	_, err := s.Client.PutObject(ctx, &s3.PutObjectInput{
		Bucket: aws.String(bucketName),
		Key:    aws.String(key),
		Body:   bytes.NewReader(content)})

	if err != nil {
		log.Printf("Failed to upload to S3: %v", err)
		return err
	}

	return nil
}
