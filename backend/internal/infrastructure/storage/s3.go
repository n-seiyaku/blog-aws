package storage

import (
	"context"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/s3"
)

type S3Storage struct {
	client    *s3.Client
	presigner *s3.PresignClient
	bucket    string
}

func NewS3Storage(
	client *s3.Client,
	presigner *s3.PresignClient,
	bucket string,
) *S3Storage {
	return &S3Storage{
		client:    client,
		presigner: presigner,
		bucket:    bucket,
	}
}

func (s *S3Storage) PresignUpload(
	ctx context.Context,
	key string,
	contentType string,
) (string, error) {
	request, err := s.presigner.PresignPutObject(
		ctx,
		&s3.PutObjectInput{
			Bucket:      aws.String(s.bucket),
			Key:         aws.String(key),
			ContentType: aws.String(contentType),
		},
		func(opts *s3.PresignOptions) {
			opts.Expires = 5 * time.Minute
		})

	if err != nil {
		return "", err
	}

	return request.URL, nil
}

func (s *S3Storage) PresignDownload(
	ctx context.Context,
	key string,
) (string, error) {
	request, err := s.presigner.PresignGetObject(
		ctx,
		&s3.GetObjectInput{
			Bucket: aws.String(s.bucket),
			Key:    aws.String(key),
		},
		func(opts *s3.PresignOptions) {
			opts.Expires = 5 * time.Minute
		},
	)
	if err != nil {
		return "", err
	}

	return request.URL, nil
}

func (s *S3Storage) Head(
	ctx context.Context,
	key string,
) (*s3.HeadObjectOutput, error) {
	return s.client.HeadObject(ctx, &s3.HeadObjectInput{
		Bucket: aws.String(s.bucket),
		Key:    aws.String(key),
	})
}

func (s *S3Storage) Delete(ctx context.Context, key string) error {
	_, err := s.client.DeleteObject(ctx, &s3.DeleteObjectInput{
		Bucket: aws.String(s.bucket),
		Key:    aws.String(key),
	})
	return err
}
