package document

import (
	"bytes"
	"context"
	"encoding/base64"
	"image"
	"image/jpeg"
	"log"
)

type document struct {
}

type DocumentService interface {
	ConvertToJpeg(ctx context.Context, imageString string) ([]byte, error)
}

func New() (DocumentService, error) {
	return &document{}, nil
}

func (d *document) ConvertToJpeg(ctx context.Context, imageString string) ([]byte, error) {
	select {
	case <-ctx.Done():
		return nil, ctx.Err()
	default:
	}

	imageData, err := base64.StdEncoding.DecodeString(imageString)
	if err != nil {
		log.Fatalf("Failed to decode string to byte image: %v", err)
	}

	img, _, err := image.Decode(bytes.NewReader(imageData))
	if err != nil {
		log.Fatalf("Failed to decode image: %v", err)
	}

	var jpegBuffer bytes.Buffer
	err = jpeg.Encode(&jpegBuffer, img, nil)
	if err != nil {
		log.Fatalf("Failed to encode image to JPEG: %v", err)
	}

	log.Printf("Successfully processed image to JPEG")

	return jpegBuffer.Bytes(), nil
}
