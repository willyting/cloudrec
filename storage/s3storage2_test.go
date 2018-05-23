package storage

import (
	"GaChaMachine/mock_s3_sdk"
	"bytes"
	"io/ioutil"
	"strings"
	"testing"

	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/aws/aws-sdk-go/service/s3/s3iface"
	"github.com/golang/mock/gomock"
)

//go:generate $GOPATH/bin/mockgen -destination src/GaChaMachine/mocks/mock_s3_client.go -package mocks github.com/aws/aws-sdk-go/service/s3/s3iface S3API

type mockS3Connecter struct {
	MockClient *mock_s3_sdk.MockS3API
}

func (s *mockS3Connecter) Connect(file *FileInfo) (s3iface.S3API, error) {
	return s.MockClient, nil
}

func TestDownload(t *testing.T) {
	testFilename := "test/test.txt"
	expected := "test"
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	mockClient := mock_s3_sdk.NewMockS3API(mockCtrl)
	connector := &mockS3Connecter{MockClient: mockClient}
	testDownloader := &S3Client{Connecter: connector}
	writer := new(bytes.Buffer)
	mockClient.EXPECT().GetObject(&s3.GetObjectInput{
		Key: &testFilename,
	}).Return(&s3.GetObjectOutput{Body: ioutil.NopCloser(strings.NewReader(expected))}, nil)
	err := testDownloader.Download(&FileInfo{
		FileName: "test/test.txt",
	}, writer)
	if err != nil {
		t.Errorf("api return error: %s", err.Error())
	}
	if writer.String() != expected {
		t.Errorf("file contant fail: got %s, want %s", writer.String(), expected)
	}
}

func TestUpload(t *testing.T) {
	var err error
	testFilename := "test/test.txt"
	expected := "test"
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	mockClient := mock_s3_sdk.NewMockS3API(mockCtrl)
	connector := &mockS3Connecter{MockClient: mockClient}
	testClient := &S3Client{Connecter: connector}
	fileContant := strings.NewReader(expected)
	mockClient.EXPECT().PutObject(gomock.Any()).Return(&s3.PutObjectOutput{}, nil).
		Do(func(in *s3.PutObjectInput) {
			if *(in.Key) != testFilename {
				t.Errorf("intput file name errror, got %s, want %s", *(in.Key), testFilename)
			}
			contant, err := ioutil.ReadAll(in.Body)
			if err != nil {
				t.Errorf("read buffer fail : %s", err.Error())
			}
			if string(contant) != expected {
				t.Errorf("put object api, got %s want %s", contant, expected)
			}
		})
	err = testClient.Upload(&FileInfo{
		FileName: "test/test.txt",
	}, fileContant)
	if err != nil {
		t.Errorf("api return error: %s", err.Error())
	}
}
