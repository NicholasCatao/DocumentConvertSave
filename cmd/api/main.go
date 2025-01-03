package main

import (
	"context"
	"encoding/json"
	"log"

	contract "mynewgolamb/domain/contracts"
	documentService "mynewgolamb/domain/document"
	uploadService "mynewgolamb/domain/upload"

	"github.com/aws/aws-lambda-go/lambda"

)

const defaultImageFormat = ".jpeg"
const defaultDirName = "Documents/"


func handleRequest(ctx context.Context, event json.RawMessage) error {
	// Parse the input event
	var document contract.Document
	if err := json.Unmarshal(event, &document); err != nil {
		log.Printf("Failed to unmarshal event: %v", err)
		return err

	}

	docService, err := documentService.New()
	if err != nil {
		log.Fatalf("Failed to initiate document service : %v", err)
	}

	jgpImage, err := docService.ConvertToJpeg(ctx, document.Image)
	if err != nil {
		log.Fatalf("Failed to decode base64 string: %v", err)
	}

	bucketName := "mylambucketgo"
	receiptContent := jgpImage
	key := defaultDirName + document.ID + defaultImageFormat

	
	uploadService, err := uploadService.New()

	if err != nil {
		log.Fatalf("Failed to initiate upload service : %v", err)
	}

	if err := uploadService.UploadReceiptToS3(ctx, bucketName, key, receiptContent); err != nil {
		return err
	}

	log.Printf("Successfully processed order %s and stored receipt in S3 bucket %s", document.ID, bucketName)

	return nil
}


func main() {
	lambda.Start(handleRequest)
}
