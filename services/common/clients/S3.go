package clients

import (
	"context"
	"io"
	"strings"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/feature/s3/manager"
	"github.com/aws/aws-sdk-go-v2/service/s3"
)

type IS3Client interface {
	ListObjects(ctx context.Context) ([]string, error)
	GetObject(ctx context.Context, key string) (io.Reader, error)
	PutObjectInput(ctx context.Context, f io.Reader, path string) error
}

var _ IS3Client = s3Client{}

type s3Client struct {
	awsS3Client *s3.Client
	bucket      string
}

// NewS3Client if no host is provided it by default create a client that connects to AWS Cloud
func NewS3Client(host, region, bucket, accessKey, pwdKey string) (IS3Client, error) {
	// FIXME(JPR): aws.EndpointResolverFunc is a deprecated method, we should use EndpointResolverWithOptionsFunc
	staticResolver := aws.EndpointResolverFunc(func(service, region string) (aws.Endpoint, error) { //nolint
		return aws.Endpoint{
			PartitionID:       "aws",
			URL:               host,
			SigningRegion:     region,
			HostnameImmutable: true,
		}, nil
	})

	// FIXME(JPR): aws.EndpointResolverFunc is a deprecated method, we should use EndpointResolverWithOptionsFunc
	cfg := aws.Config{
		Region:           region,
		Credentials:      credentials.NewStaticCredentialsProvider(accessKey, pwdKey, ""),
		EndpointResolver: staticResolver,
	}

	return &s3Client{
		awsS3Client: s3.NewFromConfig(cfg),
		bucket:      bucket,
	}, nil
}

func (s s3Client) ListObjects(ctx context.Context) ([]string, error) {
	delimiter := "/"
	input := &s3.ListObjectsV2Input{
		Bucket:    aws.String(s.bucket),
		MaxKeys:   10,
		Delimiter: &delimiter,
	}
	res, err := s.awsS3Client.ListObjectsV2(ctx, input)
	if err != nil {
		return nil, err
	}

	objectsName := make([]string, len(res.CommonPrefixes))

	for index, obj := range res.CommonPrefixes {
		if obj.Prefix == nil {
			continue
		}
		objectsName[index] = strings.TrimSuffix(*obj.Prefix, "/")
	}

	return objectsName, nil
}

func (s s3Client) GetObject(ctx context.Context, key string) (io.Reader, error) {
	input := &s3.GetObjectInput{
		Bucket: aws.String(s.bucket),
		Key:    aws.String(key),
	}

	response, err := s.awsS3Client.GetObject(ctx, input)
	if err != nil {
		return nil, err
	}

	return response.Body, nil
}

func (s s3Client) PutObjectInput(ctx context.Context, fileReader io.Reader, path string) error {
	uploader := manager.NewUploader(s.awsS3Client)

	_, err := uploader.Upload(ctx, &s3.PutObjectInput{
		Bucket: aws.String(s.bucket),
		Key:    aws.String(path),
		Body:   fileReader,
	})
	return err
}
