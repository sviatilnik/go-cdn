package storage

import (
	"bytes"
	"context"
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"io"
	"mime"
	"path/filepath"
	"strings"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/google/uuid"
	appConfig "github.com/sviatilnik/go-cdn/internal/config"
)

type S3Storage struct {
	client *s3.Client
}

func NewS3Storage(ctx context.Context, cnf *appConfig.StorageConfig) (*S3Storage, error) {
	cfg, err := config.LoadDefaultConfig(ctx,
		config.WithRegion(cnf.Region),
	)
	if err != nil {
		return nil, err
	}

	client := s3.NewFromConfig(cfg, func(o *s3.Options) {
		o.BaseEndpoint = aws.String(cnf.URL)
	})

	return &S3Storage{
		client: client,
	}, nil
}

func (s *S3Storage) SaveFile(ctx context.Context, file io.Reader, filename string) (*File, error) {
	bucketName := s.getBucketName()
	if err := s.ensureBucketExists(ctx, bucketName); err != nil {
		return nil, err
	}

	ext := filepath.Ext(filename)
	newFileName := fmt.Sprintf("%s%s", uuid.New().String(), ext)

	_, err := s.client.PutObject(ctx, &s3.PutObjectInput{
		Bucket:      aws.String(bucketName),
		Key:         aws.String(newFileName),
		Body:        file,
		ContentType: aws.String(mime.TypeByExtension(filepath.Ext(filename))),
	})

	if err != nil {
		return nil, err
	}

	return &File{
		Filename:    newFileName,
		Path:        fmt.Sprintf("%s/%s", bucketName, newFileName),
		Size:        0,
		ContentType: mime.TypeByExtension(ext),
		Timestamp:   time.Now(),
	}, nil
}

func (s *S3Storage) GetFile(ctx context.Context, relativePath string) (*File, error) {
	bucketName, key := s.getBucketNameAndKey(relativePath)

	result, err := s.client.GetObject(ctx, &s3.GetObjectInput{
		Bucket: aws.String(bucketName),
		Key:    aws.String(key),
	})
	if err != nil {
		if strings.Contains(err.Error(), "NoSuchKey") {
			return nil, ErrFileNotFound
		}
		return nil, err
	}

	defer result.Body.Close()

	body, err := io.ReadAll(result.Body)
	if err != nil {
		return nil, err
	}

	return &File{
		Filename:    key,
		Path:        relativePath,
		Size:        int64(len(body)),
		ContentType: mime.TypeByExtension(filepath.Ext(key)),
		Timestamp:   time.Now(),
		File:        &readSeekCloser{bytes.NewReader(body)},
	}, nil
}

func (s *S3Storage) DeleteFile(ctx context.Context, relativePath string) error {
	bucketName, key := s.getBucketNameAndKey(relativePath)

	_, err := s.client.DeleteObject(ctx, &s3.DeleteObjectInput{
		Bucket: aws.String(bucketName),
		Key:    aws.String(key),
	})
	if err != nil {
		return err
	}

	return nil
}

func (s *S3Storage) GetFileMaxSize() int64 {
	return 10 << 20
}

func (s *S3Storage) ensureBucketExists(ctx context.Context, bucketName string) error {
	result, err := s.client.ListBuckets(ctx, &s3.ListBucketsInput{})
	if err != nil {
		return err
	}

	for _, bucket := range result.Buckets {
		if *bucket.Name == bucketName {
			return nil
		}
	}
	_, err = s.client.CreateBucket(ctx, &s3.CreateBucketInput{
		Bucket: aws.String(bucketName),
	})
	if err != nil {
		return err
	}

	return nil
}

func (s *S3Storage) getBucketName() string {
	hash := hex.EncodeToString(md5.New().Sum([]byte(time.Now().Format("200601021504"))))
	hash = string(hash[len(hash)-5:])

	return hash
}

func (s *S3Storage) getBucketNameAndKey(relPath string) (bucketName, key string) {
	if strings.Contains(relPath, "/") {
		path := strings.Split(relPath, "/")

		bucketName = path[0]
		key = path[1]
	} else {
		bucketName = s.getBucketName()
		key = relPath
	}

	return
}

type readSeekCloser struct {
	*bytes.Reader
}

func (r *readSeekCloser) Close() error {
	return nil
}
