package filestore

import (
	"mime/multipart"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
)

type FileStore interface {
	UploadFile(file multipart.File, fileKey string) (string, error)
	Delete(fileKey string) error
}

type MockFileStore struct{}

func (m *MockFileStore) UploadFile(file multipart.File, name string) (string, error) {
  return "http://example.com/" + name, nil
}

func (m *MockFileStore) Delete(fileURL string) error {
  return nil
}

type AWSFileStore struct {
	Bucket string
	S3Sess *session.Session
}

func (fs *AWSFileStore) UploadFile(file multipart.File, s3Key string) (string, error) {
	uploader := s3manager.NewUploader(fs.S3Sess)

	result, err := uploader.Upload(&s3manager.UploadInput{
		Bucket: aws.String(fs.Bucket),
		Key:    aws.String(s3Key),
		Body:   file,
	})
	if err != nil {
		return "", err
	}

	return result.Location, nil
}

func (fs *AWSFileStore) Delete(s3Key string) error {
	s3Client := s3.New(fs.S3Sess)
	_, err := s3Client.DeleteObject(&s3.DeleteObjectInput{
		Bucket: aws.String(fs.Bucket),
		Key:    aws.String(s3Key),
	})
	if err != nil {
		return err
	}
	return nil
}
