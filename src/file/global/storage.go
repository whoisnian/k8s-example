package global

import (
	"context"
	"errors"
	"io"

	"github.com/whoisnian/k8s-example/src/file/pkg/fsdriver"
	"github.com/whoisnian/k8s-example/src/file/pkg/s3driver"
)

var FS StorageDriver

type StorageDriver interface {
	SetupTracing()

	CreateFile(ctx context.Context, bucket, object string, reader io.Reader, size int64) (int64, error)
	OpenFile(ctx context.Context, bucket, object string) (io.ReadCloser, error)
	DeleteFile(ctx context.Context, bucket, object string) error
}

func SetupStorage() {
	var err error
	if CFG.StorageDriver == "filesystem" {
		FS, err = fsdriver.New(CFG.RootDirectory)
	} else if CFG.StorageDriver == "aws-s3" {
		FS, err = s3driver.New(CFG.S3Endpoint, CFG.S3AccessKey, CFG.S3SecretKey, CFG.S3Secure)
	} else {
		err = errors.New("unknown storage driver")
	}
	if err != nil {
		panic(err)
	}
	FS.SetupTracing()
}
