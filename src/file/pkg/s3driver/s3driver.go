package s3driver

import (
	"context"
	"errors"
	"io"
	"io/fs"

	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/trace"
)

type Driver struct {
	client *minio.Client
	tracer trace.Tracer
}

func New(endpoint, accessKey, secretKey string, secure bool) (*Driver, error) {
	client, err := minio.New(endpoint, &minio.Options{
		Creds:  credentials.NewStaticV4(accessKey, secretKey, ""),
		Secure: secure,
	})
	if err != nil {
		return nil, err
	}
	return &Driver{client: client}, nil
}

func (dri *Driver) SetupTracing() {
	if dri.tracer == nil {
		dri.tracer = otel.GetTracerProvider().Tracer("github.com/whoisnian/k8s-example/src/file/pkg/s3driver")
	}
}

func (dri *Driver) CreateFile(ctx context.Context, bucket, object string, reader io.Reader, size int64) (int64, error) {
	ctx, span := dri.tracer.Start(ctx, "storage.CreateFile")
	defer span.End()

	info, err := dri.createFile(ctx, bucket, object, reader, size)
	if err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())
	}
	return info.Size, err
}

func (dri *Driver) createFile(ctx context.Context, bucket, object string, reader io.Reader, size int64) (info minio.UploadInfo, err error) {
	if !fs.ValidPath(bucket) || !fs.ValidPath(object) {
		return info, errors.New("s3driver: invalid bucket/object to resolve")
	}

	if ok, err := dri.client.BucketExists(ctx, bucket); err != nil {
		return info, err
	} else if !ok {
		if err = dri.client.MakeBucket(ctx, bucket, minio.MakeBucketOptions{}); err != nil {
			return info, err
		}
	}
	return dri.client.PutObject(ctx, bucket, object, reader, size, minio.PutObjectOptions{})
}

func (dri *Driver) OpenFile(ctx context.Context, bucket, object string) (io.ReadCloser, error) {
	ctx, span := dri.tracer.Start(ctx, "storage.OpenFile")
	defer span.End()

	obj, err := dri.client.GetObject(ctx, bucket, object, minio.GetObjectOptions{})
	if err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())
	}
	return obj, err
}

func (dri *Driver) DeleteFile(ctx context.Context, bucket, object string) error {
	ctx, span := dri.tracer.Start(ctx, "storage.DeleteFile")
	defer span.End()

	err := dri.client.RemoveObject(ctx, bucket, object, minio.RemoveObjectOptions{})
	if err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())
	}
	return err
}
