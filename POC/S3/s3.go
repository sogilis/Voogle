package main

import (
	"context"
	"fmt"
	"log"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/service/s3"
)

func main() {
	creds := credentials.NewStaticCredentialsProvider("key", "secret", "")

	cfg, err := config.LoadDefaultConfig(context.TODO(), config.WithCredentialsProvider(creds), config.WithRegion("eu-west-3"))
	if err != nil {
		log.Printf("error: %v", err)
		return
	}

	awsS3Client := s3.NewFromConfig(cfg)

	input := &s3.ListObjectsInput{
		Bucket:  aws.String("voogle-video"),
		MaxKeys: 10,
	}
	res, err := awsS3Client.ListObjects(context.Background(), input)
	if err != nil {
		fmt.Println("ERR", err)
		return
	}

	for _, obj := range res.Contents {
		// Do whatever you need with each object "obj"
		fmt.Println("ETAG", *obj.Key)
	}
}
