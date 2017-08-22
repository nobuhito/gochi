package gochi

import (
	"context"
	"io/ioutil"

	"cloud.google.com/go/storage"
	"google.golang.org/appengine/file"
)

type Storage struct {
	Path string
	Attr *storage.ObjectAttrs
	Body []byte
}

func (s *Storage) Get(ctx context.Context) error {
	bucketname, err := file.DefaultBucketName(ctx)
	if err != nil {
		return err
	}

	client, err := storage.NewClient(ctx)
	if err != nil {
		return err
	}
	defer client.Close()

	obj := client.Bucket(bucketname).Object(s.Path)
	s.Attr, err = obj.Attrs(ctx)
	if err != nil {
		return err
	}

	r, err := obj.NewReader(ctx)
	if err != nil {
		return err
	}
	defer r.Close()

	body, err := ioutil.ReadAll(r)
	if err != nil {
		return err
	}
	s.Body = body

	return nil
}

func (s *Storage) Write(ctx context.Context) error {
	bucketname, err := file.DefaultBucketName(ctx)
	if err != nil {
		return err
	}

	client, err := storage.NewClient(ctx)
	if err != nil {
		return err
	}
	defer client.Close()

	w := client.Bucket(bucketname).Object(s.Path).NewWriter(ctx)
	_, err = w.Write(s.Body)
	if err != nil {
		return err
	}
	defer w.Close()

	return nil
}

func (s *Storage) Delete(ctx context.Context) error {
	bucketname, err := file.DefaultBucketName(ctx)
	if err != nil {
		return err
	}

	client, err := storage.NewClient(ctx)
	if err != nil {
		return err
	}
	defer client.Close()

	err = client.Bucket(bucketname).Object(s.Path).Delete(ctx)
	if err != nil {
		return err
	}

	return nil
}
