package s3driver

import (
	"context"
	"io"

	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
)

type Driver struct {
	client      *minio.Client
	bucketCache map[string]bool
}

func New(endpoint, accessKey, secretKey string, secure bool) (*Driver, error) {
	client, err := minio.New(endpoint, &minio.Options{
		Creds:  credentials.NewStaticV4(accessKey, secretKey, ""),
		Secure: secure,
	})
	if err != nil {
		return nil, err
	}
	return &Driver{
		client:      client,
		bucketCache: make(map[string]bool),
	}, nil
}

func (dri *Driver) CreateFile(bucket, object string, reader io.Reader, size int64) (string, int64, error) {
	ctx := context.Background()
	if err := dri.ensureBucketExists(ctx, bucket); err != nil {
		return "", 0, err
	}

	info, err := dri.client.PutObject(ctx, bucket, object, reader, size, minio.PutObjectOptions{})
	if err != nil {
		return "", 0, err
	}
	return info.ChecksumSHA256, info.Size, nil
}

func (dri *Driver) OpenFile(bucket, object string) (io.ReadCloser, error) {
	return dri.client.GetObject(context.Background(), bucket, object, minio.GetObjectOptions{})
}

func (dri *Driver) DeleteFile(bucket, object string) error {
	return dri.client.RemoveObject(context.Background(), bucket, object, minio.RemoveObjectOptions{})
}

func (dri *Driver) ensureBucketExists(ctx context.Context, bucket string) error {
	if !dri.bucketCache[bucket] {
		if ok, err := dri.client.BucketExists(ctx, bucket); err != nil {
			return err
		} else if !ok {
			if err = dri.client.MakeBucket(ctx, bucket, minio.MakeBucketOptions{}); err != nil {
				return err
			}
		}
		dri.bucketCache[bucket] = true
	}
	return nil
}
