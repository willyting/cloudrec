package storage

import (
	"bytes"
	"io"
	"io/ioutil"
	"path/filepath"

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
func (s *S3Stroage) GetDownloader() Downloader {
	return &S3Client{Connecter: &S3Connecter{}}
}

// GetUploader ...
func (s *S3Stroage) GetUploader() Uploader {
	return &S3Client{Connecter: &S3Connecter{}}
}

// GetLister ...
func (s *S3Stroage) GetLister() Lister {
	return &S3Client{Connecter: &S3Connecter{}}
}

// Download ...
func (s *S3Client) Download(file *FileInfo, writeTo io.Writer) error {
	client, err := s.Connect(file)
	if err != nil {
		return err
	}
	out, err := client.GetObject(&s3.GetObjectInput{
		Key: aws.String(file.FileName),
	})
	if err != nil {
		return err
	}
	_, err = io.Copy(writeTo, out.Body)
	out.Body.Close()
	if err != nil {
		return err
	}
	return nil
}

// Upload ...
func (s *S3Client) Upload(file *FileInfo, readFrom io.ReadCloser) error {
	client, err := s.Connect(file)
	if err != nil {
		return err
	}
	buffer, err := ioutil.ReadAll(readFrom)
	readFrom.Close()
	if err != nil {
		return err
	}
	_, err = client.PutObject(&s3.PutObjectInput{
		Bucket: aws.String(file.Bucket),
		Body:   io.ReadSeeker(bytes.NewReader(buffer)),
		Key:    aws.String(file.FileName)})
	if err != nil {
		return err
	}
	return nil
}

// List ...
func (s *S3Client) List(file *FileInfo) ([]string, error) {
	client, err := s.Connect(file)
	if err != nil {
		return nil, err
	}
	out, err := client.ListObjectsV2(&s3.ListObjectsV2Input{
		Bucket: aws.String(file.Bucket),
		Prefix: aws.String(file.FileName),
	})
	if err != nil {
		return nil, err
	}
	var list []string
	for _, obj := range out.Contents {
		if filepath.Ext(*obj.Key) == ".db" {
			list = append(list, *obj.Key)
		}
	}
	return list, nil
}

// Connect return a storage client support all s3 API
// return client on success, and return nil and an error on fail
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
