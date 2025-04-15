package asset_store

import (
	"bytes"
	"fmt"
	"strings"
	"time"

	"github.com/PAY-HERO-CONSULTING/gh-tools/logger"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
)

type assetStoreConfig struct {
	Bucket          string
	Region          string
	AccessKeyID     string
	SecretAccessKey string
	Endpoint        string
	Environment     string
}

func NewAssetStoreConfig(
	bucket,
	region,
	accessKeyID,
	secretAccessKey,
	endpoint,
	environment string,
) *assetStoreConfig {
	return &assetStoreConfig{
		Bucket:          bucket,
		Region:          region,
		AccessKeyID:     accessKeyID,
		SecretAccessKey: secretAccessKey,
		Endpoint:        endpoint,
		Environment:     environment,
	}
}

type AssetStore interface {
	Delete(key string) error
	DownloadURL(key string) (string, error)
	Upload(b []byte, contentType, key string, public bool) (string, error)
}

type S3AssetStore struct {
	development bool
	s3          *s3.S3
	bucket      string
	region      string
	endpoint    string
}

func NewAssetStore(
	config *assetStoreConfig,
) *S3AssetStore {
	return NewAssetStoreWithCredentials(
		config.Bucket,
		config.Region,
		config.AccessKeyID,
		config.SecretAccessKey,
		config.Endpoint,
		config.Environment,
	)
}

func NewAssetStoreWithCredentials(
	bucket,
	region,
	accessKeyId,
	secretAccessKey,
	endpoint,
	environment string,
) *S3AssetStore {
	config := &aws.Config{
		Region: aws.String(region),
		Credentials: credentials.NewStaticCredentials(
			accessKeyId,
			secretAccessKey,
			"",
		),
	}

	if endpoint != "" {
		config.Endpoint = aws.String(endpoint)
		config.S3ForcePathStyle = aws.Bool(true)
	}

	s3Session, err := session.NewSession()
	if err != nil {
		panic("failed to initiate AWS S3 session")
	}

	awsS3 := s3.New(s3Session, config)

	return &S3AssetStore{
		development: strings.ToLower(environment) == "development",
		s3:          awsS3,
		bucket:      bucket,
		region:      region,
		endpoint:    endpoint,
	}
}

func (s *S3AssetStore) Delete(key string) error {
	input := &s3.DeleteObjectInput{
		Bucket: aws.String(s.bucket),
		Key:    aws.String(key),
	}

	output, err := s.s3.DeleteObject(input)
	if err != nil {
		return err
	}

	logger.Infof("Successfully deleted file bucket: %v, key: %v, output: %v", s.bucket, key, output.GoString())

	return nil
}

func (s *S3AssetStore) DownloadURL(key string) (string, error) {
	req, _ := s.s3.GetObjectRequest(
		&s3.GetObjectInput{
			Bucket: aws.String(s.bucket),
			Key:    aws.String(key),
		},
	)

	downloadURL, err := req.Presign(60 * time.Minute)
	if err != nil {
		return downloadURL, err
	}

	if s.development {
		downloadURL = strings.Replace(downloadURL, "http://localstack", "http://localhost", -1)
	}

	return downloadURL, nil
}

func (s *S3AssetStore) Upload(b []byte, contentType, key string, public bool) (string, error) {
	acl := "private"
	if public {
		acl = "public-read"
	}

	body := bytes.NewReader(b)

	params := &s3.PutObjectInput{
		ACL:           aws.String(acl),
		Body:          body,
		Bucket:        aws.String(s.bucket),
		ContentLength: aws.Int64(body.Size()),
		Key:           aws.String(key),
	}

	if contentType != "" {
		params.ContentType = aws.String(contentType)
	}

	result, err := s.s3.PutObject(params)
	if err != nil {
		return "", err
	}

	logger.Infof(
		"[AssetStorage Upload] Successfully uploaded file bucket: %v, key: %v, contentType: %v, size: %v, etag: %v",
		aws.StringValue(params.Bucket),
		aws.StringValue(params.Key),
		aws.StringValue(params.ContentType),
		aws.Int64Value(params.ContentLength),
		aws.StringValue(result.ETag),
	)

	// Construct the file URL
	url := ""
	if public {
		url = fmt.Sprintf("https://%s.s3.%s.amazonaws.com%s", s.bucket, s.region, key)
	} else {
		// For private files, generate a presigned URL
		req, _ := s.s3.GetObjectRequest(&s3.GetObjectInput{
			Bucket: aws.String(s.bucket),
			Key:    aws.String(key),
		})
		url, err = req.Presign(15 * time.Minute) // Valid for 15 minutes
		if err != nil {
			return "", fmt.Errorf("failed to generate presigned URL: %w", err)
		}
	}

	return url, nil
}
