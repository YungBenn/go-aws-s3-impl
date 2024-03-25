package main

import (
	"bytes"
	"context"
	"encoding/base64"
	"strings"

	"log"
	"net/http"
	"os"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/feature/s3/manager"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/aws/aws-sdk-go-v2/service/s3/types"
)

var base64Data = ""

func main() {
	// uploadTom()
	uploadBase64()
}

func uploadBase64() {
	// Load the Shared AWS Configuration (~/.aws/config)
	cfg, err := config.LoadDefaultConfig(context.TODO(),
		config.WithSharedCredentialsFiles(
			[]string{"aws/credentials"},
		),
		config.WithSharedConfigFiles(
			[]string{"aws/config"},
		),
	)
	if err != nil {
		log.Fatalf("unable to load SDK config, %v", err)
	}

	// Create an Amazon S3 service client
	client := s3.NewFromConfig(cfg)

	// Create an uploader with the client and default options
	uploader := manager.NewUploader(client)

	// Check if the base64 data string contains a comma
	if strings.Contains(base64Data, ",") {
		// Remove the prefix
		base64Data = strings.Split(base64Data, ",")[1]
	} else {
		log.Fatalf("invalid base64 data")
	}

	// Decode the base64 data
	data, err := base64.StdEncoding.DecodeString(base64Data)
	if err != nil {
		log.Fatalf("unable to decode base64 data, %v", err)
	}

	// Convert the data to a reader
	reader := bytes.NewReader(data)

	// Upload the file to S3
	result, err := uploader.Upload(context.TODO(), &s3.PutObjectInput{
		Bucket:      aws.String("yungbenn"),
		Key:         aws.String("tom-base64"),
		Body:        reader,
		ACL:         types.ObjectCannedACLPublicRead,
		ContentType: aws.String("image/jpeg"),
	})

	if err != nil {
		log.Fatalf("unable to upload file, %v", err)
	}

	log.Println(result)
}

func uploadTom() {
	// Load the Shared AWS Configuration (~/.aws/config)
	cfg, err := config.LoadDefaultConfig(context.TODO(),
		config.WithSharedCredentialsFiles(
			[]string{"aws/credentials"},
		),
		config.WithSharedConfigFiles(
			[]string{"aws/config"},
		),
	)
	if err != nil {
		log.Fatalf("unable to load SDK config, %v", err)
	}

	// Open the file
	file, err := os.Open("tom.jpg")
	if err != nil {
		log.Fatalf("unable to open file, %v", err)
	}
	defer file.Close()

	// Read the first 512 bytes of the file to determine its content type
	buffer := make([]byte, 512)
	_, err = file.Read(buffer)
	if err != nil {
		log.Fatalf("unable to read file, %v", err)
	}

	// Reset the file pointer back to the start of the file
	file.Seek(0, 0)

	// Determine the content type of the file
	contentType := http.DetectContentType(buffer)

	// Create an Amazon S3 service client
	client := s3.NewFromConfig(cfg)

	uploader := manager.NewUploader(client)
	result, err := uploader.Upload(context.TODO(), &s3.PutObjectInput{
		Bucket:      aws.String("yungbenn"),
		Key:         aws.String("tom"),
		Body:        file,
		ACL:         types.ObjectCannedACLPublicRead,
		ContentType: aws.String(contentType),
	})

	if err != nil {
		log.Fatalf("unable to upload file, %v", err)
	}

	log.Println(result)
}
