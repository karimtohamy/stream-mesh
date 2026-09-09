package storage

import (
	"context"
	"fmt"
	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
	"path/filepath"
	"stream-mesh/media-sync/internal/config"
	"strings"
	"os"
)

type Storage interface {
	Exists(ctx context.Context, bucketName, key string) (bool, error)
	Download(ctx context.Context, bucketName, key, localPath string) error
	UploadDirectory(ctx context.Context, bucketName, localPath, keyPrefix string) error
	EnsureBucket(ctx context.Context, bucketName string) error
}

type MinIOStorage struct {
	client *minio.Client
}

func NewMinIOStorage(cfg *config.MinIOConfig) (*MinIOStorage, error) {
	client, err := minio.New(cfg.Endpoint, &minio.Options{
		Creds:  credentials.NewStaticV4(cfg.AccessKey, cfg.SecretKey, ""),
		Secure: cfg.UseSSL,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to init client: %w", err)
	}
	return &MinIOStorage{
		client: client,
	}, nil
}

func (m *MinIOStorage) EnsureBucket(ctx context.Context, bucketName string) error {
	exists, err := m.client.BucketExists(ctx, bucketName)
	if err != nil {
		return fmt.Errorf("failed to check bucket %s: %w", bucketName, err)
	}
	if !exists {
		err = m.client.MakeBucket(ctx, bucketName, minio.MakeBucketOptions{})
		if err != nil {
			return fmt.Errorf("failed to create bucket %s: %w", bucketName, err)
		}
	}
	return nil
}

func (m *MinIOStorage) Exists(ctx context.Context, bucketName, key string) (bool, error) {
	_, err := m.client.StatObject(ctx, bucketName, key, minio.StatObjectOptions{})
	if err != nil {
		errResp := minio.ToErrorResponse(err)
		if errResp.Code == "NoSuchKey" || errResp.Code == "NoSuchBucket" {
			return false, nil
		}
		return false, fmt.Errorf("failed to stat object: %w", err)
	}
	return true, nil
}

func (m *MinIOStorage) Download(ctx context.Context, bucketName, key, localPath string) error {
	return m.client.FGetObject(ctx, bucketName, key, localPath, minio.GetObjectOptions{})
}

func (m *MinIOStorage) UploadDirectory(ctx context.Context, bucketName, localPath, keyPrefix string) error {
	if err := m.EnsureBucket(ctx, bucketName); err != nil {
		return err
	}
	return filepath.Walk(localPath, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if info.IsDir() {
			return nil
		}

		relPath, err := filepath.Rel(localPath, path)
		if err != nil {
			return fmt.Errorf("failed to calculate relative path: %w", err)
		}

		objectKey := filepath.ToSlash(filepath.Join(keyPrefix, relPath))

		contentType := "application/octet-stream"
		switch {
		case strings.HasSuffix(path, ".m3u8"):
			contentType = "application/vnd.apple.mpegurl"
		case strings.HasSuffix(path, ".ts"):
			contentType = "video/MP2T"
		}

		_, err = m.client.FPutObject(ctx, bucketName, objectKey, path, minio.PutObjectOptions{
			ContentType: contentType,
		})
		if err != nil {
			return fmt.Errorf("failed to upload %s: %w", objectKey, err)
		}

		return nil
	})
}
