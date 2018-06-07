package storage

import (
	"bytes"
	"gachamachine/mock_s3_sdk"
	"io"
	"io/ioutil"
	"testing"

	"github.com/golang/mock/gomock"
)

type mockFileSystem struct {
	MockClient *mock_s3_sdk.MockReadWriteCloser
}

func (m *mockFileSystem) Open(file *FileInfo) (io.ReadWriteCloser, error) {
	return m.MockClient, nil
}

func TestFileOperator_Download(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	mockClient := mock_s3_sdk.NewMockReadWriteCloser(mockCtrl)
	connector := &mockFileSystem{MockClient: mockClient}
	testDownloader := &FileOperator{SystemOperator: connector}
	mockClient.EXPECT().Close().Return(nil)
	type args struct {
		file *FileInfo
	}
	tests := []struct {
		name        string
		s           *FileOperator
		args        args
		wantWriteTo string
		wantErr     bool
	}{
		// TODO: Add test cases.
		{"normal", testDownloader, args{&FileInfo{FileName: "test/test.txt"}}, "test", false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			getRead := mockClient.EXPECT().Read(gomock.Any()).Times(1).DoAndReturn(func(p []byte) (n int, err error) {
				if p == nil {
					t.Error("buffer error")
				}
				p[0] = 't'
				p[1] = 'e'
				p[2] = 's'
				p[3] = 't'
				return 4, nil
			})
			mockClient.EXPECT().Read(gomock.Any()).After(getRead).Return(0, io.EOF)
			writeTo := &bytes.Buffer{}
			if err := tt.s.Download(tt.args.file, writeTo); (err != nil) != tt.wantErr {
				t.Errorf("FileOperator.Download() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if gotWriteTo := writeTo.String(); gotWriteTo != tt.wantWriteTo {
				t.Errorf("FileOperator.Download() = %v, want %v", gotWriteTo, tt.wantWriteTo)
			}
		})
	}
}

func TestFileOperator_Upload(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	mockClient := mock_s3_sdk.NewMockReadWriteCloser(mockCtrl)
	connector := &mockFileSystem{MockClient: mockClient}
	testDownloader := &FileOperator{SystemOperator: connector}
	mockClient.EXPECT().Close().Return(nil)
	type args struct {
		file     *FileInfo
		readFrom io.ReadCloser
	}
	tests := []struct {
		name            string
		s               *FileOperator
		args            args
		wantWriteToMock string
		wantErr         bool
	}{
		// TODO: Add test cases.
		{"normal", testDownloader,
			args{&FileInfo{FileName: "test/test.txt"}, ioutil.NopCloser(bytes.NewReader([]byte("test")))},
			"test", false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockClient.EXPECT().Write([]byte(tt.wantWriteToMock)).Return(len(tt.wantWriteToMock), nil)
			if err := tt.s.Upload(tt.args.file, tt.args.readFrom); (err != nil) != tt.wantErr {
				t.Errorf("FileOperator.Upload() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
