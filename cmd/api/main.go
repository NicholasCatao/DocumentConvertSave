package main

import (
	"bytes"
	"context"
	"encoding/base64"
	"encoding/json"
	"image"
	"image/jpeg"
	"log"

	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/s3"
)

const defaultImageFormat = ".jpeg"
const defaultDirName = "Documents/"

type Document struct {
	ID    string `json:"id"`
	Image string `json:"image"`
}

var (
	s3Client *s3.Client
)

func init() {
	// Initialize the S3 client outside of the handler, during the init phase
	cfg, err := config.LoadDefaultConfig(context.TODO())
	if err != nil {
		log.Fatalf("unable to load SDK config, %v", err)
	}

	s3Client = s3.NewFromConfig(cfg)
}

func uploadReceiptToS3(ctx context.Context, bucketName string, key string, receiptContent []byte) error {
	_, err := s3Client.PutObject(ctx, &s3.PutObjectInput{
		Bucket: &bucketName,
		Key:    &key,
		Body:   bytes.NewReader(receiptContent),
	})

	if err != nil {
		log.Printf("Failed to upload receipt to S3: %v", err)
		return err
	}
	return nil
}

func handleRequest(ctx context.Context, event json.RawMessage) error {
	// Parse the input event
	var document Document
	if err := json.Unmarshal(event, &document); err != nil {
		log.Printf("Failed to unmarshal event: %v", err)
		return err

	}

	jgpImage, err := ConvertToJpeg(document.Image)

	if err != nil {
		log.Fatalf("Failed to decode base64 string: %v", err)
	}

	bucketName := "mylambucketgo"
	receiptContent := jgpImage
	key := defaultDirName + document.ID + defaultImageFormat

	// Upload the receipt to S3 using the helper method
	if err := uploadReceiptToS3(ctx, bucketName, key, receiptContent); err != nil {
		return err
	}

	log.Printf("Successfully processed order %s and stored receipt in S3 bucket %s", document.ID, bucketName)

	return nil
}

func ConvertToJpeg(imageString string) ([]byte, error) {

	imageData, err := base64.StdEncoding.DecodeString(imageString)

	if err != nil {
		log.Fatalf("Failed to decode string to byte image: %v", err)
	}

	image, _, err := image.Decode(bytes.NewReader(imageData))

	if err != nil {
		log.Fatalf("Failed to decode image: %v", err)
	}

	var jpegBuffer bytes.Buffer

	err = jpeg.Encode(&jpegBuffer, image, nil)

	if err != nil {
		log.Fatalf("Failed to encode image to JPEG: %v", err)
	}

	if err != nil {
		return nil, err
	}

	log.Printf("Successfully processed image to JPEG")

	return jpegBuffer.Bytes(), nil
}

func main() {
	lambda.Start(handleRequest)
}
