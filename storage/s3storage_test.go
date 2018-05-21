package storage_test

import (
	"GaChaMachine/mocks"
	"GaChaMachine/storage"
	"bytes"
	"io/ioutil"
	"strings"
	"testing"

	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/aws/aws-sdk-go/service/s3/s3iface"
	"github.com/golang/mock/gomock"
)

type mockS3Connecter struct {
	MockClient *mocks.MockS3API
}

func (s *mockS3Connecter) Connect(file *storage.FileInfo) (s3iface.S3API, error) {
	return s.MockClient, nil
}

func TestDownload(t *testing.T) {
	testFilename := "test/test.txt"
	expected := "test"
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	mockClient := mocks.NewMockS3API(mockCtrl)
	connector := &mockS3Connecter{MockClient: mockClient}
	testDownloader := &storage.S3Client{Connecter: connector}
	writer := new(bytes.Buffer)
	mockClient.EXPECT().GetObject(&s3.GetObjectInput{
		Key: &testFilename,
	}).Return(&s3.GetObjectOutput{Body: ioutil.NopCloser(strings.NewReader(expected))}, nil)
	err := testDownloader.Download(&storage.FileInfo{
		FileName: "test/test.txt",
	}, writer)
	if err != nil {
		t.Errorf("api return error: %s", err.Error())
	}
	if writer.String() != expected {
		t.Errorf("file contant fail: got %s, want %s", writer.String(), expected)
	}
}
