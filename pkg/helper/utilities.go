package helpers

import (
	"fmt"
	"io/ioutil"
	"mime/multipart"
	"os"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
)

func IsLambda() bool {
	if lambdaTaskRoot := os.Getenv("LAMBDA_TASK_ROOT"); lambdaTaskRoot != "" {
		return true
	}
	return false
}

func UploadToS3(file *multipart.FileHeader, filePath, awsRegion, s3Bucket string) error {
	fileBytes, err := file.Open()
	if err != nil {
		fmt.Printf("Unable to upload due to error : %v", err)
		return err
	}
	defer fileBytes.Close()

	// Create an AWS session
	sess, err := session.NewSession(&aws.Config{
		Region: aws.String(awsRegion),
	})
	if err != nil {
		fmt.Printf("Unable to upload due to error : %v", err)
		return err
	}

	// Create an S3 service client
	s3Client := s3.New(sess)

	// Prepare the S3 input parameters
	params := &s3.PutObjectInput{
		Bucket: aws.String(s3Bucket),
		Key:    aws.String(filePath),
		Body:   fileBytes,
		ACL:    aws.String("public-read"), // Adjust permissions as needed
	}

	// Upload the file to S3
	_, err = s3Client.PutObject(params)
	if err != nil {
		fmt.Printf("Unable to upload due to error : %v", err)
		return err
	}

	return nil
}

func DownloadFromS3(filePath, awsRegion, s3Bucket string) ([]byte, error) {
	// Create an AWS session
	sess, err := session.NewSession(&aws.Config{
		Region: aws.String(awsRegion),
	})
	if err != nil {
		fmt.Printf("Unable to download due to error: %v", err)
		return nil, err
	}

	// Create an S3 service client
	s3Client := s3.New(sess)

	// Prepare the S3 input parameters
	params := &s3.GetObjectInput{
		Bucket: aws.String(s3Bucket),
		Key:    aws.String(filePath),
	}

	// Download the file from S3
	resp, err := s3Client.GetObject(params)
	if err != nil {
		fmt.Printf("Unable to download due to error: %v", err)
		return nil, err
	}
	defer resp.Body.Close()

	// Read the S3 object content into a byte slice
	imageBytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Printf("Unable to read content due to error: %v", err)
		return nil, err
	}

	fmt.Println("Download successful!")
	return imageBytes, nil
}
