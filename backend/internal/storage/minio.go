package storage

import (
	"context"
	"io"

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
	_, err := m.Client.PutObject(
		ctx,
		m.Bucket,
		objectKey,
		reader,
		size,
		minio.PutObjectOptions{
			ContentType: contentType,
		},
	)

	return err
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
