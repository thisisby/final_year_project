package s3

import (
	"context"
	"fmt"
	"mime/multipart"
	"path/filepath"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/s3"
)

type Client struct {
	client *s3.Client
}

func NewClient(region string, accessKeyID string, secretAccessKey string) (*Client, error) {
	cfg, err := config.LoadDefaultConfig(
		context.TODO(),
		config.WithRegion(region),
		config.WithCredentialsProvider(aws.CredentialsProviderFunc(func(ctx context.Context) (aws.Credentials, error) {
			return aws.Credentials{
				AccessKeyID:     accessKeyID,
				SecretAccessKey: secretAccessKey,
			}, nil
		})),
	)
	if err != nil {
		return nil, err
	}

	s3Client := s3.NewFromConfig(cfg)

	return &Client{client: s3Client}, nil
}

func (c *Client) UploadFile(file multipart.File, filename string, bucket string) (string, error) {
	ctx := context.TODO() // Consider using a more specific context.

	contentType := "application/octet-stream" // Default if detection fails.
	ext := filepath.Ext(filename)
	switch ext {
	case ".jpg", ".jpeg":
		contentType = "image/jpeg"
	case ".png":
		contentType = "image/png"
	case ".gif":
		contentType = "image/gif"
	}

	fmt.Println("Bucket name:", bucket)
	fmt.Println("File name:", filename)
	fmt.Println("Content Type:", contentType)

	_, err := c.client.PutObject(ctx, &s3.PutObjectInput{
		Bucket:      aws.String(bucket),
		Key:         aws.String(filename), // Or generate a unique key.
		Body:        file,
		ContentType: &contentType,
	})
	if err != nil {
		fmt.Println("Error uploading file:", err)
		return "", err
	}

	// Construct the URL.
	url := "https://" + bucket + ".s3.us-west-2.amazonaws.com/" + filename

	return url, nil
}
