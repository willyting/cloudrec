package storage

import (
	"gachamachine/mock_s3_sdk"
	"testing"

	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/aws/aws-sdk-go/service/s3/s3iface"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
)

func TestS3Client_List(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	testList := &S3Client{Connecter: &mockS3Connecter{MockClient: createMockS3List(mockCtrl)}}
	type args struct {
		file *FileInfo
	}
	tests := []struct {
		name    string
		s       *S3Client
		args    args
		want    []string
		wantErr bool
	}{
		// TODO: Add test cases.
		{"one-file", testList, args{&FileInfo{FileName: "user/test/2018-05-11"}},
			[]string{"user/test/2018-05-11/1/t.db"}, false},
		{"many-file-one-db", testList, args{&FileInfo{FileName: "user/test/2018-05-09"}},
			[]string{"user/test/2018-05-09/1/t.db"}, false},
		{"many-file-one-db-first", testList, args{&FileInfo{FileName: "user/test/2018-05-13"}},
			[]string{"user/test/2018-05-13/1/t.db"}, false},
		{"many-file-one-db-last", testList, args{&FileInfo{FileName: "user/test/2018-05-10"}},
			[]string{"user/test/2018-05-10/1/t.db"}, false},
		{"many-file-two-db", testList, args{&FileInfo{FileName: "user/test/2018-05-08"}},
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
func createMockS3List(ctrl *gomock.Controller) *mock_s3_sdk.MockS3API {
	mock := mock_s3_sdk.NewMockS3API(ctrl)
	dates := []struct {
		prefix string   // input: filename
		list   []string // output: database list
	}{
		{"user/test/2018-05-06", []string{
			"user/test/2018-05-06/1/t.db",
			"user/test/2018-05-06/2/t.db",
			"user/test/2018-05-06/3/t.db"}},
		{"user/test/2018-05-08", []string{
			"user/test/2018-05-08/1/1.3gp",
			"user/test/2018-05-08/1/2.3gp",
			"user/test/2018-05-08/1/t.db",
			"user/test/2018-05-08/1/3.3gp",
			"user/test/2018-05-08/2/t.db",
			"user/test/2018-05-08/2/1.3gp",
			"user/test/2018-05-08/2/2.3gp"}},
		{"user/test/2018-05-09", []string{
			"user/test/2018-05-09/1/1.3gp",
			"user/test/2018-05-09/1/2.3gp",
			"user/test/2018-05-09/1/t.db",
			"user/test/2018-05-09/1/3.3gp"}},
		{"user/test/2018-05-10", []string{
			"user/test/2018-05-10/1/1.3gp",
			"user/test/2018-05-10/1/2.3gp",
			"user/test/2018-05-10/1/3.3gp",
			"user/test/2018-05-10/1/t.db"}},
		{"user/test/2018-05-11", []string{
			"user/test/2018-05-11/1/t.db"}},
		{"user/test/2018-05-13", []string{
			"user/test/2018-05-13/1/t.db",
			"user/test/2018-05-13/1/1.3gp",
			"user/test/2018-05-13/1/2.3gp",
			"user/test/2018-05-13/1/3.3gp"}},
	}
	mock.EXPECT().ListObjectsV2(gomock.Any()).
		DoAndReturn(func(input *s3.ListObjectsV2Input) (*s3.ListObjectsV2Output, error) {
			for _, dd := range dates {
				if *input.Prefix == dd.prefix {
					out := new(s3.ListObjectsV2Output)
					out.Contents = make([]*s3.Object, len(dd.list))
					for i := range dd.list {
						out.Contents[i] = new(s3.Object)
						out.Contents[i].Key = new(string)
						*(out.Contents[i].Key) = dd.list[i]
					}
					return out, nil
				}
			}
			return nil, nil
		}).AnyTimes()
	return mock
}

func TestS3Connecter_Connect(t *testing.T) {
	type args struct {
		file *FileInfo
	}
	tests := []struct {
		name    string
		s       *S3Connecter
		args    args
		want    s3iface.S3API
		wantErr bool
	}{
	// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.s.Connect(tt.args.file)
			if assert.EqualValues(t, tt.wantErr, err != nil) {
				return
			}
			assert.EqualValues(t, tt.want, got)
		})
	}
}
