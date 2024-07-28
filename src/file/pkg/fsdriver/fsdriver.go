package fsdriver

import (
	"context"
	"errors"
	"io"
	"io/fs"
	"os"
	"path/filepath"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/trace"
)

type Driver struct {
	root   string
	tracer trace.Tracer
}

func New(root string) (*Driver, error) {
	if root == "" {
		return nil, errors.New("fsdriver: empty root directory name")
	}
	if info, err := os.Stat(root); err != nil {
		return nil, err
	} else if !info.IsDir() {
		return nil, errors.New("fsdriver: root is not a directory")
	}
	return &Driver{root: root}, nil
}

func (dri *Driver) SetupTracing() {
	if dri.tracer == nil {
		dri.tracer = otel.GetTracerProvider().Tracer("github.com/whoisnian/k8s-example/src/file/pkg/fsdriver")
	}
}

func (dri *Driver) CreateFile(ctx context.Context, bucket, object string, reader io.Reader, _ int64) (int64, error) {
	ctx, span := dri.tracer.Start(ctx, "storage.CreateFile")
	defer span.End()

	written, err := dri.createFile(ctx, bucket, object, reader)
	if err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())
	}
	return written, err
}

func (dri *Driver) createFile(_ context.Context, bucket, object string, reader io.Reader) (int64, error) {
	name, err := dri.resolve(bucket, object)
	if err != nil {
		return 0, err
	}

	parent := filepath.Dir(name)
	if info, err := os.Stat(parent); err != nil && os.IsNotExist(err) {
		if err = os.MkdirAll(parent, 0755); err != nil {
			return 0, err
		}
	} else if err != nil {
		return 0, err
	} else if !info.IsDir() {
		return 0, errors.New("fsdriver: parent is not a directory")
	}

	file, err := os.Create(name)
	if err != nil {
		return 0, err
	}
	defer file.Close()

	return io.Copy(file, reader)
}

func (dri *Driver) OpenFile(ctx context.Context, bucket, object string) (io.ReadCloser, error) {
	ctx, span := dri.tracer.Start(ctx, "storage.OpenFile")
	defer span.End()

	fi, err := dri.openFile(ctx, bucket, object)
	if err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())
	}
	return fi, err
}

func (dri *Driver) openFile(_ context.Context, bucket, object string) (*os.File, error) {
	name, err := dri.resolve(bucket, object)
	if err != nil {
		return nil, err
	}
	return os.Open(name)
}

func (dri *Driver) DeleteFile(ctx context.Context, bucket, object string) error {
	ctx, span := dri.tracer.Start(ctx, "storage.DeleteFile")
	defer span.End()

	err := dri.deleteFile(ctx, bucket, object)
	if err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())
	}
	return err
}

func (dri *Driver) deleteFile(_ context.Context, bucket, object string) error {
	name, err := dri.resolve(bucket, object)
	if err != nil {
		return err
	}
	return os.Remove(name)
}

func (dri *Driver) resolve(bucket, object string) (string, error) {
	if !fs.ValidPath(bucket) || !fs.ValidPath(object) {
		return "", errors.New("fsdriver: invalid bucket/object to resolve")
	}
	return filepath.Join(dri.root, bucket, object), nil
}
