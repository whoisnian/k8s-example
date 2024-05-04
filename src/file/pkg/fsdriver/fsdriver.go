package fsdriver

import (
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"io"
	"io/fs"
	"os"
	"path/filepath"
)

type Driver struct {
	root string
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
	return &Driver{root}, nil
}

func (dri *Driver) CreateFile(bucket, object string, reader io.Reader, _ int64) (string, int64, error) {
	name, err := dri.resolve(bucket, object)
	if err != nil {
		return "", 0, err
	}

	parent := filepath.Dir(name)
	if info, err := os.Stat(parent); err != nil && os.IsNotExist(err) {
		if err = os.MkdirAll(parent, 0755); err != nil {
			return "", 0, err
		}
	} else if err != nil {
		return "", 0, err
	} else if !info.IsDir() {
		return "", 0, errors.New("fsdriver: parent is not a directory")
	}

	file, err := os.Create(name)
	if err != nil {
		return "", 0, err
	}
	defer file.Close()

	hasher := sha256.New()
	size, err := io.Copy(io.MultiWriter(file, hasher), reader)
	if err != nil {
		return "", 0, err
	}
	return hex.EncodeToString(hasher.Sum(nil)), size, nil
}

func (dri *Driver) OpenFile(bucket, object string) (io.ReadCloser, error) {
	name, err := dri.resolve(bucket, object)
	if err != nil {
		return nil, err
	}
	return os.Open(name)
}

func (dri *Driver) DeleteFile(bucket, object string) error {
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