package store

import (
	"context"
	"fmt"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/service/s3"
)

type S3Config struct {
	Client       *s3.Client
	BucketName   string
	BaseEndpoint string
}

// NewS3 creates a new S3Config with Cloudflare R2 storage configuration
// Parameters:
//
//	bucketName - Name of the R2 bucket to use
//	accountId - Cloudflare account ID
//	accessKeyId - R2 access key ID
//	accessKeySecret - R2 access key secret
//
// Returns:
//
//	S3Config - Configuration for R2 storage operations
//	error - Any error that occurred during setup
func NewS3(bucketName, accountId, accessKeyId, accessKeySecret string) (S3Config, error) {
	if bucketName == "" || accountId == "" || accessKeyId == "" || accessKeySecret == "" {
		return S3Config{}, fmt.Errorf("s3 parameters are required")
	}

	cfg, err := config.LoadDefaultConfig(
		context.TODO(),
		config.WithCredentialsProvider(
			credentials.NewStaticCredentialsProvider(accessKeyId, accessKeySecret, "")),
		config.WithRegion("auto"),
	)
	if err != nil {
		return S3Config{}, fmt.Errorf("unable to load S3 config: %v", err)
	}

	baseEndpoint := fmt.Sprintf("https://%s.r2.cloudflarestorage.com", accountId)
	client := s3.NewFromConfig(cfg, func(o *s3.Options) {
		o.BaseEndpoint = aws.String(baseEndpoint)
	})

	return S3Config{
		Client:       client,
		BucketName:   bucketName,
		BaseEndpoint: baseEndpoint,
	}, nil
}
