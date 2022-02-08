package clients

import (
	"context"
	"io"
)

var _ IS3Client = s3ClientDummy{}

type s3ClientDummy struct {
	listObjects    func() ([]string, error)
	getObject      func(id string) (io.Reader, error)
	putObjectInput func(f io.Reader, title string) error
}

func NewS3ClientDummy(listObjects func() ([]string, error), getObject func(string) (io.Reader, error), putObjectInput func(io.Reader, string) error) IS3Client {
	return s3ClientDummy{listObjects, getObject, putObjectInput}
}

func (s s3ClientDummy) ListObjects(ctx context.Context) ([]string, error) {
	return s.listObjects()
}

func (s s3ClientDummy) GetObject(ctx context.Context, id string) (io.Reader, error) {
	return s.getObject(id)
}

func (s s3ClientDummy) PutObjectInput(ctx context.Context, f io.Reader, title string) error {
	return s.putObjectInput(f, title)
}
