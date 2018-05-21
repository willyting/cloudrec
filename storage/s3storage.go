package storage

import (
	"fmt"
	"io"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/aws/aws-sdk-go/service/s3/s3iface"
)

// Connecter ...
type Connecter interface {
	Connect(*FileInfo) (s3iface.S3API, error)
}

// S3Stroage ...
type S3Stroage struct {
}

// S3Client ...
type S3Client struct {
	Connecter
}

// S3Connecter ...
type S3Connecter struct {
}

// GetDownloader ...
func (s *S3Stroage) GetDownloader() *S3Client {
	return &S3Client{Connecter: &S3Connecter{}}
}

// GetUploader ...
func (s *S3Stroage) GetUploader() *S3Client {
	return &S3Client{Connecter: &S3Connecter{}}
}

// Download ...
func (s *S3Client) Download(file *FileInfo, writeTo io.Writer) error {
	client, err := s.Connect(file)
	if err != nil {
		return err
	}
	if client == nil {
		return fmt.Errorf("create client fail")
	}
	out, _ := client.GetObject(&s3.GetObjectInput{
		Key: aws.String(file.FileName),
	})
	_, err = io.Copy(writeTo, out.Body)
	out.Body.Close()
	if err != nil {
		return err
	}
	return nil
}

// Connect ...
func (s *S3Connecter) Connect(file *FileInfo) (s3iface.S3API, error) {
	cred := credentials.NewStaticCredentials(file.AccessKeyID, file.SecretKey, file.SessionToken)
	cred.Get()
	config := aws.NewConfig().
		WithCredentials(cred).
		WithRegion(file.Region)
	sess, err := session.NewSession(config)
	if err != nil {
		return nil, err
	}
	return s3.New(sess), nil
}
