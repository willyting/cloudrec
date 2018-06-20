package storage

//go:generate %GOPATH%\bin\mockgen -destination src\gachamachine\mock_s3_sdk\mock_local_dir.go -package mock_s3_sdk gachamachine/storage/iface FolderOperator
import (
	"bytes"
	"gachamachine/mock_s3_sdk"
	"gachamachine/storage/iface"
	"io"
	"io/ioutil"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
)

type mockFileSystem struct {
	MockClient *mock_s3_sdk.MockReadWriteCloser
	MockDir    *mock_s3_sdk.MockFolderOperator
}

func (m *mockFileSystem) Open(file *FileInfo) (io.ReadWriteCloser, error) {
	return m.MockClient, nil
}
func (m *mockFileSystem) OpenDir(file *FileInfo) (iface.FolderOperator, error) {
	return m.MockDir, nil
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
				assert.NotNil(t, p, "buffer error")
				p[0] = 't'
				p[1] = 'e'
				p[2] = 's'
				p[3] = 't'
				return 4, nil
			})
			mockClient.EXPECT().Read(gomock.Any()).After(getRead).Return(0, io.EOF)
			writeTo := &bytes.Buffer{}
			err := tt.s.Download(tt.args.file, writeTo)
			if assert.EqualValues(t, tt.wantErr, err != nil) {
				return
			}
			assert.EqualValues(t, tt.wantWriteTo, writeTo.String())
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
			err := tt.s.Upload(tt.args.file, tt.args.readFrom)
			assert.EqualValues(t, tt.wantErr, err != nil)
		})
	}
}

func TestFileOperator_List(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	type args struct {
		file *FileInfo
	}
	tests := []struct {
		name    string
		s       *FileOperator
		args    args
		want    []string
		wantErr bool
	}{
		// TODO: Add test cases.
		{"normal", &FileOperator{SystemOperator: &mockFileSystem{MockDir: createMockDir(mockCtrl)}},
			args{&FileInfo{FileName: "user/test/2018-05-11"}},
			[]string{"user/test/2018-05-11/1/t.db"}, false},
		{"many-file-one-db", &FileOperator{SystemOperator: &mockFileSystem{MockDir: createMockDir(mockCtrl)}},
			args{&FileInfo{FileName: "user/test/2018-05-09"}},
			[]string{"user/test/2018-05-09/1/t.db"}, false},
		{"many-file-one-db-first", &FileOperator{SystemOperator: &mockFileSystem{MockDir: createMockDir(mockCtrl)}},
			args{&FileInfo{FileName: "user/test/2018-05-13"}},
			[]string{"user/test/2018-05-13/1/t.db"}, false},
		{"many-file-one-db-last", &FileOperator{SystemOperator: &mockFileSystem{MockDir: createMockDir(mockCtrl)}},
			args{&FileInfo{FileName: "user/test/2018-05-10"}},
			[]string{"user/test/2018-05-10/1/t.db"}, false},
		{"many-file-two-db", &FileOperator{SystemOperator: &mockFileSystem{MockDir: createMockDir(mockCtrl)}},
			args{&FileInfo{FileName: "user/test/2018-05-08"}},
			[]string{"user/test/2018-05-08/1/t.db",
				"user/test/2018-05-08/2/t.db"}, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.s.List(tt.args.file)
			if !assert.EqualValues(t, tt.wantErr, (err != nil)) {
				return
			}
			assert.EqualValues(t, tt.want, got)
		})
	}
}
func createMockDir(ctrl *gomock.Controller) *mock_s3_sdk.MockFolderOperator {
	mock := mock_s3_sdk.NewMockFolderOperator(ctrl)
	mock.EXPECT().Readdirnames(gomock.Any()).Times(1).Return([]string{
		"user.test.2018-05-06.1.t.db",
		"user.test.2018-05-06.2.t.db",
		"user.test.2018-05-06.3.t.db",
		"user.test.2018-05-08.1.1.3gp",
		"user.test.2018-05-08.1.2.3gp",
		"user.test.2018-05-08.1.t.db",
		"user.test.2018-05-08.1.3.3gp",
		"user.test.2018-05-08.2.t.db",
		"user.test.2018-05-08.2.1.3gp",
		"user.test.2018-05-08.2.2.3gp",
		"user.test.2018-05-09.1.1.3gp",
		"user.test.2018-05-09.1.2.3gp",
		"user.test.2018-05-09.1.t.db",
		"user.test.2018-05-09.1.3.3gp",
		"user.test.2018-05-10.1.1.3gp",
		"user.test.2018-05-10.1.2.3gp",
		"user.test.2018-05-10.1.3.3gp",
		"user.test.2018-05-10.1.t.db",
		"user.test.2018-05-11.1.t.db",
		"user.test.2018-05-13.1.t.db",
		"user.test.2018-05-13.1.1.3gp",
		"user.test.2018-05-13.1.2.3gp",
		"user.test.2018-05-13.1.3.3gp"}, nil)
	return mock
}
