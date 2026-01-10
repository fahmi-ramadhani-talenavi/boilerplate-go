// Package storage provides cloud storage utilities for file operations.
// Supports AWS S3 and S3-compatible services like MinIO.
package storage

import (
	"context"
	"fmt"
	"io"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/user/go-boilerplate/pkg/logger"
	"go.uber.org/zap"
)

// S3Config holds S3 client configuration
type S3Config struct {
	Region    string
	Bucket    string
	AccessKey string
	SecretKey string
	Endpoint  string // Optional: for MinIO or LocalStack
}

// S3Client wraps the AWS S3 client with helper methods
type S3Client struct {
	client       *s3.Client
	presignClient *s3.PresignClient
	bucket       string
}

// NewS3Client creates a new S3 client instance
func NewS3Client(cfg S3Config) (*S3Client, error) {
	// Create AWS config with static credentials
	awsCfg, err := config.LoadDefaultConfig(context.Background(),
		config.WithRegion(cfg.Region),
		config.WithCredentialsProvider(credentials.NewStaticCredentialsProvider(
			cfg.AccessKey,
			cfg.SecretKey,
			"",
		)),
	)
	if err != nil {
		return nil, fmt.Errorf("failed to load AWS config: %w", err)
	}

	// Create S3 client options
	opts := []func(*s3.Options){}
	
	// Use custom endpoint for MinIO or LocalStack
	if cfg.Endpoint != "" {
		opts = append(opts, func(o *s3.Options) {
			o.BaseEndpoint = aws.String(cfg.Endpoint)
			o.UsePathStyle = true // Required for MinIO
		})
	}

	client := s3.NewFromConfig(awsCfg, opts...)
	presignClient := s3.NewPresignClient(client)

	return &S3Client{
		client:       client,
		presignClient: presignClient,
		bucket:       cfg.Bucket,
	}, nil
}

// Upload uploads a file to S3
//
// Parameters:
// - ctx: Context for cancellation
// - key: Object key (file path in bucket)
// - reader: File content reader
// - contentType: MIME type of the file
//
// Returns:
// - string: URL of uploaded file
// - error: If upload fails
func (s *S3Client) Upload(ctx context.Context, key string, reader io.Reader, contentType string) (string, error) {
	input := &s3.PutObjectInput{
		Bucket:      aws.String(s.bucket),
		Key:         aws.String(key),
		Body:        reader,
		ContentType: aws.String(contentType),
	}

	_, err := s.client.PutObject(ctx, input)
	if err != nil {
		logger.Error(ctx, "Failed to upload to S3", zap.Error(err), zap.String("key", key))
		return "", fmt.Errorf("failed to upload file: %w", err)
	}

	logger.Info(ctx, "File uploaded to S3", zap.String("key", key))
	return fmt.Sprintf("s3://%s/%s", s.bucket, key), nil
}

// Download retrieves a file from S3
//
// Parameters:
// - ctx: Context for cancellation
// - key: Object key to download
//
// Returns:
// - io.ReadCloser: File content reader (caller must close)
// - error: If download fails
func (s *S3Client) Download(ctx context.Context, key string) (io.ReadCloser, error) {
	input := &s3.GetObjectInput{
		Bucket: aws.String(s.bucket),
		Key:    aws.String(key),
	}

	result, err := s.client.GetObject(ctx, input)
	if err != nil {
		logger.Error(ctx, "Failed to download from S3", zap.Error(err), zap.String("key", key))
		return nil, fmt.Errorf("failed to download file: %w", err)
	}

	return result.Body, nil
}

// Delete removes a file from S3
//
// Parameters:
// - ctx: Context for cancellation
// - key: Object key to delete
//
// Returns:
// - error: If deletion fails
func (s *S3Client) Delete(ctx context.Context, key string) error {
	input := &s3.DeleteObjectInput{
		Bucket: aws.String(s.bucket),
		Key:    aws.String(key),
	}

	_, err := s.client.DeleteObject(ctx, input)
	if err != nil {
		logger.Error(ctx, "Failed to delete from S3", zap.Error(err), zap.String("key", key))
		return fmt.Errorf("failed to delete file: %w", err)
	}

	logger.Info(ctx, "File deleted from S3", zap.String("key", key))
	return nil
}

// GetPresignedURL generates a presigned URL for temporary access
//
// Parameters:
// - ctx: Context for cancellation
// - key: Object key
// - expiry: How long the URL should be valid
//
// Returns:
// - string: Presigned URL
// - error: If generation fails
func (s *S3Client) GetPresignedURL(ctx context.Context, key string, expiry time.Duration) (string, error) {
	input := &s3.GetObjectInput{
		Bucket: aws.String(s.bucket),
		Key:    aws.String(key),
	}

	presignedReq, err := s.presignClient.PresignGetObject(ctx, input, func(opts *s3.PresignOptions) {
		opts.Expires = expiry
	})
	if err != nil {
		logger.Error(ctx, "Failed to generate presigned URL", zap.Error(err), zap.String("key", key))
		return "", fmt.Errorf("failed to generate presigned URL: %w", err)
	}

	return presignedReq.URL, nil
}

// GetPresignedUploadURL generates a presigned URL for uploading
//
// Parameters:
// - ctx: Context for cancellation
// - key: Object key where file will be uploaded
// - contentType: Expected MIME type
// - expiry: How long the URL should be valid
//
// Returns:
// - string: Presigned upload URL
// - error: If generation fails
func (s *S3Client) GetPresignedUploadURL(ctx context.Context, key string, contentType string, expiry time.Duration) (string, error) {
	input := &s3.PutObjectInput{
		Bucket:      aws.String(s.bucket),
		Key:         aws.String(key),
		ContentType: aws.String(contentType),
	}

	presignedReq, err := s.presignClient.PresignPutObject(ctx, input, func(opts *s3.PresignOptions) {
		opts.Expires = expiry
	})
	if err != nil {
		logger.Error(ctx, "Failed to generate presigned upload URL", zap.Error(err), zap.String("key", key))
		return "", fmt.Errorf("failed to generate presigned upload URL: %w", err)
	}

	return presignedReq.URL, nil
}
