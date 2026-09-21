package storage

import (
	"context"
	"fmt"
	"io"
	"log"

	"github.com/minio/minio-go/v7"
)

type ImageStorage interface {
	Upload(
		ctx context.Context,
		objectKey string,
		reader io.Reader,
		size int64,
		contentType string,
	) error

	Delete(
		ctx context.Context,
		objectKey string,
	) error
}

type MinioStorage struct {
	Client *minio.Client
	Bucket string
}

func (m *MinioStorage) Upload(
	ctx context.Context,
	objectKey string,
	reader io.Reader,
	size int64,
	contentType string,
) error {
	log.Printf("[MINIO DEBUG] Attempting upload -> Bucket: %q, Key: %q, Size: %d, ContentType: %q",
		m.Bucket, objectKey, size, contentType)

	info, err := m.Client.PutObject(
		ctx,
		m.Bucket,
		objectKey,
		reader,
		size,
		minio.PutObjectOptions{
			ContentType: contentType,
		},
	)

	if err != nil {
		log.Printf("[MINIO DEBUG] Upload FAILED -> Error: %v", err)
		return fmt.Errorf("minio put object: %w", err)
	}

	log.Printf("[MINIO DEBUG] Upload SUCCESS -> Key: %q, Uploaded Size: %d", info.Key, info.Size)
	return nil
}

func (m *MinioStorage) Delete(
	ctx context.Context,
	objectKey string,
) error {
	return m.Client.RemoveObject(
		ctx,
		m.Bucket,
		objectKey,
		minio.RemoveObjectOptions{},
	)
}
