package helpers

import (
	"fmt"
	"io/ioutil"
	"mime/multipart"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/jinzhu/gorm"
	"github.com/saiteja111997/throttle_backend/pkg/structures"
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

func DeleteFromS3(objectKey, awsRegion, s3Bucket string) error {
	// Create an AWS session
	sess, err := session.NewSession(&aws.Config{
		Region: aws.String(awsRegion),
	})
	if err != nil {
		fmt.Printf("Unable to create session for deleting from S3: %v\n", err)
		return err
	}

	// Create an S3 service client
	s3Client := s3.New(sess)

	// Prepare the S3 input parameters
	params := &s3.DeleteObjectInput{
		Bucket: aws.String(s3Bucket),
		Key:    aws.String(objectKey),
	}

	// Delete the object from S3
	_, err = s3Client.DeleteObject(params)
	if err != nil {
		fmt.Printf("Unable to delete from S3 due to error: %v\n", err)
		return err
	}

	// Wait until the object is deleted to ensure it's gone
	err = s3Client.WaitUntilObjectNotExists(&s3.HeadObjectInput{
		Bucket: aws.String(s3Bucket),
		Key:    aws.String(objectKey),
	})
	if err != nil {
		fmt.Printf("Error occurred while waiting for the object to be deleted from S3: %v\n", err)
		return err
	}

	fmt.Println("Successfully deleted the document from S3!")
	return nil
}

func HashPassword(password string) (string, error) {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	return string(hashedPassword), nil
}

func UpdateDocFilePath(db *gorm.DB, errorID string, newDocFilePath string) error {
	result := db.Model(&structures.Errors{}).Where("id = ?", errorID).Update("doc_file_path", newDocFilePath)
	return result.Error
}

func UpdateDocStatus(db *gorm.DB, errorID, status string) error {

	statusInt, err := strconv.Atoi(status)
	if err != nil {
		// Handle error, e.g., log or return
		fmt.Printf("Error converting status to int: %v\n", err)
		return err
	}

	result := db.Model(&structures.Errors{}).Where("id = ?", errorID).Update("status", statusInt)
	return result.Error

}

func DeleteDocFromDB(db *gorm.DB, errorID string) error {
	result := db.Where("id = ?", errorID).Delete(&structures.Errors{})
	if result.Error != nil {
		fmt.Printf("Unable to delete doc from DB: %v", result.Error)
		return result.Error
	}
	return nil
}

// Helper function to format date with suffixes (st, nd, rd, th)
func FormatDate(t time.Time) string {
	// Create the base date format
	// formattedDate := t.Format("Jan 2, 2006")

	// Get the day of the month
	day := t.Day()

	// Determine the suffix for the day (st, nd, rd, th)
	var suffix string
	if day%10 == 1 && day != 11 {
		suffix = "st"
	} else if day%10 == 2 && day != 12 {
		suffix = "nd"
	} else if day%10 == 3 && day != 13 {
		suffix = "rd"
	} else {
		suffix = "th"
	}

	// Insert the suffix into the formatted date
	return t.Format("Jan 2") + suffix + t.Format(", 2006")
}
