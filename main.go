package main

import (
	"context"
	"fmt"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"log"
)

const (
	region     = "ap-east-1"
	bucketName = "hanzi-read-test"
	maxKeys    = 0

	accessKeyId     = "AKIAZBI2PYWUKHGJY5MC"
	accessSecretKey = "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"
)

func main() {

	// Load the SDK's configuration from environment and shared config, and
	// create the client with this.
	cfg, err := config.LoadDefaultConfig(context.TODO(), config.WithRegion(region), config.WithCredentialsProvider(credentials.StaticCredentialsProvider{
		Value: aws.Credentials{
			AccessKeyID: accessKeyId, SecretAccessKey: accessSecretKey, SessionToken: "",
		},
	}))
	if err != nil {
		log.Fatalf("failed to load SDK configuration, %v", err)
	}

	client := s3.NewFromConfig(cfg)

	var bn = bucketName
	// Set the parameters based on the CLI flag inputs.
	params := &s3.ListObjectsV2Input{
		Bucket: &bn,
	}
	//if len(objectPrefix) != 0 {
	//	params.Prefix = &objectPrefix
	//}
	//if len(objectDelimiter) != 0 {
	//	params.Delimiter = &objectDelimiter
	//}

	// Create the Paginator for the ListObjectsV2 operation.
	p := s3.NewListObjectsV2Paginator(client, params, func(o *s3.ListObjectsV2PaginatorOptions) {
		if v := int32(maxKeys); v != 0 {
			o.Limit = v
		}
	})

	var i int
	for p.HasMorePages() {
		i++
		// Next Page takes a new context for each page retrieval. This is where
		// you could add timeouts or deadlines.
		page, err := p.NextPage(context.TODO())
		if err != nil {
			log.Fatalf("failed to get page %v, %v", i, err)
		}

		// Log the objects found
		for _, obj := range page.Contents {
			fmt.Println("Object:", *obj.Key)
		}
	}
}
