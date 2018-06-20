package storage

//go:generate $GOPATH/bin/mockgen -destination src/gachamachine/mock_s3_sdk/mock_s3_client.go -package mock_s3_sdk github.com/aws/aws-sdk-go/service/s3/s3iface S3API
import (
	"bytes"
	"gachamachine/mock_s3_sdk"
	"io/ioutil"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/aws/aws-sdk-go/service/s3/s3iface"
	"github.com/golang/mock/gomock"
)

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
	assert.NoError(t, err)
	assert.EqualValues(t, expected, writer.String())
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
	mockClient.EXPECT().PutObject(gomock.Any()).Return(&s3.PutObjectOutput{}, nil).
		Do(func(in *s3.PutObjectInput) {
			assert.EqualValues(t, testFilename, *(in.Key))
			contant, err := ioutil.ReadAll(in.Body)
			assert.NoError(t, err)
			assert.EqualValues(t, expected, string(contant))
		})
	err = testClient.Upload(&FileInfo{
		FileName: "test/test.txt",
	}, ioutil.NopCloser(strings.NewReader(expected)))
	assert.EqualValues(t, nil, err)
}
