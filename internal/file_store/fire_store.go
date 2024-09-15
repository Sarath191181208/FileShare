package filestore

import (
	"mime/multipart"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
)

type FileStore struct {
	Bucket string
	S3Sess *session.Session
}

func (fs *FileStore) UploadFile(file multipart.File, s3Key string) (string, error) {
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

func (fs *FileStore) Delete(s3Key string) error {
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
