package upload

import (
	"context"
	"log"
	externalServices "mynewgolamb/infra/externalServices"
)

type uploadService struct {
	externalServices.S3Client
}

type UploadService interface {
	UploadReceiptToS3(ctx context.Context, bucketName string, key string, receiptContent []byte) error
}

func New() (UploadService, error) {
    s3Client, err := externalServices.NewS3Client()
    if err != nil {
        return nil, err
    }
    return &uploadService{
        S3Client: s3Client,
    }, nil
}

func (u *uploadService) UploadReceiptToS3(ctx context.Context, bucketName string, key string, receiptContent []byte) error {

	err := u.UploadToS3(ctx, bucketName, key, receiptContent)

	if err != nil {
		log.Printf("Failed to upload receipt to S3: %v", err)
		return err
	}
	return nil
}