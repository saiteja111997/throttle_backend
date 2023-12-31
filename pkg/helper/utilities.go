package helpers

import (
	"fmt"
	"io/ioutil"
	"mime/multipart"
	"os"
	"strings"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"golang.org/x/crypto/bcrypt"
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

// UploadToS3 uploads text content to an S3 bucket.
func UploadTextToS3(textContent, filePath, awsRegion, s3Bucket string) error {
	// Convert the text content to a reader
	fileBytes := strings.NewReader(textContent)

	// Create an AWS session
	sess, err := session.NewSession(&aws.Config{
		Region: aws.String(awsRegion),
	})
	if err != nil {
		fmt.Printf("Unable to upload due to error: %v", err)
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

	// Upload the text content to S3
	_, err = s3Client.PutObject(params)
	if err != nil {
		fmt.Printf("Unable to upload due to error: %v", err)
		return err
	}

	return nil
}

// DownloadTextFromS3 downloads text content from an S3 bucket.
func DownloadTextFromS3(objectKey, awsRegion, s3Bucket string) (string, error) {
	// Create an AWS session
	sess, err := session.NewSession(&aws.Config{
		Region: aws.String(awsRegion),
	})
	if err != nil {
		fmt.Printf("Unable to download due to error: %v", err)
		return "", err
	}

	// Create an S3 service client
	s3Client := s3.New(sess)

	// Prepare the S3 input parameters
	params := &s3.GetObjectInput{
		Bucket: aws.String(s3Bucket),
		Key:    aws.String(objectKey),
	}

	// Download the text content from S3
	resp, err := s3Client.GetObject(params)
	if err != nil {
		fmt.Printf("Unable to download due to error: %v", err)
		return "", err
	}
	defer resp.Body.Close()

	// Read the content from the response body
	textContent, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Printf("Unable to read response body: %v", err)
		return "", err
	}

	return string(textContent), nil
}

func HashPassword(password string) (string, error) {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	return string(hashedPassword), nil
}
