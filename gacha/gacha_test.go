package gacha

import (
	"fmt"
	"gachamachine/mocks"
	"gachamachine/storage"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/gorilla/mux"
	"github.com/stretchr/testify/assert"
)

type req struct {
	url      string
	userHead string
}
type awsCred struct {
	ID    string // input: accessKeyID
	Key   string // input: secretKey
	Token string // input: filename
}

var mockCredentials = []awsCred{
	{"aaaaaaaaaaaaaaaaaaaaa",
		"bbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbb",
		""},
	{"tttttttttttttttttttt",
		"dddddddddddddddddddddddddddddddddddddddd",
		"yyyyyyyyyyyyyyyyyyyyyyyyyyyyyyyyyyyyyyyyyyyyyyyyyyyyyyyyyyyyyyyy"},
}

type recCredCase struct {
	name       string
	cerd       *awsCred
	want       string
	wantStatus int
}

func TestGetRecCrendential(t *testing.T) {
	queryReq := &req{"http://localhost/recstorage/test?p=test.txt", "user"}
	// setup mock module
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	defaultService.cloud = &mockStorage{Down: createMockCredGetRec(mockCtrl)}
	router := mux.NewRouter().StrictSlash(true)
	router.Methods("GET").Path("/recstorage/{cameraid}").Name("test").HandlerFunc(GetRec)
	tests := getGetRecCredTestCase()
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			rr := httptest.NewRecorder()
			router.ServeHTTP(rr, createReqGetRec(queryReq, tt.cerd, t))
			assert.Equal(t, tt.wantStatus, rr.Code, "wrong response status")
			assert.Equal(t, tt.want, rr.Body.String(), "wrong response body")
		})
	}
}

func createReqGetRec(r *req, c *awsCred, t *testing.T) *http.Request {
	return addCredHeader(createGetReq(r, t), c, t)
}

func getGetRecCredTestCase() []recCredCase {
	return []recCredCase{
		// TODO: Add test cases.
		{"ami-user-cred",
			&awsCred{"aaaaaaaaaaaaaaaaaaaaa",
				"bbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbb",
				""},
			"test", 200},
		{"temp-cred",
			&awsCred{"tttttttttttttttttttt",
				"dddddddddddddddddddddddddddddddddddddddd",
				"yyyyyyyyyyyyyyyyyyyyyyyyyyyyyyyyyyyyyyyyyyyyyyyyyyyyyyyyyyyyyyyy"},
			"test", 200},
		{"fail-cred",
			&awsCred{"aaaaaaaaaaaaaaaaaa",
				"bbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbb",
				""},
			"401 permission denied\n", 401},
		{"no-cred",
			&awsCred{"",
				"",
				""},
			"401 permission denied\n", 401},
		{"miss-id",
			&awsCred{"",
				"bbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbb",
				""},
			"401 permission denied\n", 401},
		{"miss-key",
			&awsCred{"tttttttttttttttttttt",
				"",
				""},
			"401 permission denied\n", 401},
	}
}

func createMockCredGetRec(ctrl *gomock.Controller) *mocks.MockDownloader {
	mock := mocks.NewMockDownloader(ctrl)
	mock.EXPECT().Download(gomock.Any(), gomock.Any()).
		DoAndReturn(func(f *storage.FileInfo, w io.Writer) error {
			for _, dd := range mockCredentials {
				if f.AccessKeyID == dd.ID &&
					f.SecretKey == dd.Key &&
					f.SessionToken == dd.Token {
					fmt.Fprint(w, "test")
					return nil
				}
			}
			return fmt.Errorf("401 permission denied")
		}).AnyTimes()
	return mock
}

func TestPutRec(t *testing.T) {
	queryReq := &req{"http://localhost/recstorage/test?p=test.txt", "user"}
	// setup mock module
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	defaultService.cloud = &mockStorage{Up: createMockCredPutRec(mockCtrl)}
	router := mux.NewRouter().StrictSlash(true)
	router.Methods("POST").Path("/recstorage/{cameraid}").Name("test").HandlerFunc(PutRec)
	tests := getGetRecCredTestCase()
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			rr := httptest.NewRecorder()
			router.ServeHTTP(rr, createReqPutRec(queryReq, tt.cerd, t))
			assert.Equal(t, tt.wantStatus, rr.Code, "wrong response status")
		})
	}
}
func createMockCredPutRec(ctrl *gomock.Controller) *mocks.MockUploader {
	mock := mocks.NewMockUploader(ctrl)
	mock.EXPECT().Upload(gomock.Any(), gomock.Any()).
		DoAndReturn(func(f *storage.FileInfo, r io.Reader) error {
			for _, dd := range mockCredentials {
				if f.AccessKeyID == dd.ID &&
					f.SecretKey == dd.Key &&
					f.SessionToken == dd.Token {
					return nil
				}
			}
			return fmt.Errorf("401 permission denied")
		}).AnyTimes()
	return mock
}
func createReqPutRec(r *req, c *awsCred, t *testing.T) *http.Request {
	return addCredHeader(createPostReq(r, t), c, t)
}

type listDBCase struct {
	name       string
	req        *req // input: info to create http req
	want       string
	wantStatus int
}

func TestListDb(t *testing.T) {
	// setup mock module
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	defaultService.cloud = &mockStorage{List: createMockList(mockCtrl)}
	router := mux.NewRouter().StrictSlash(true)
	router.Methods("GET").Path("/recstorage/{cameraid}/date").Name("testListDb").HandlerFunc(ListDb)
	tests := getListDBTestCase()
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			rr := httptest.NewRecorder()
			router.ServeHTTP(rr, createReqListDB(tt.req, t))
			assert.Equal(t, tt.wantStatus, rr.Code, "wrong response status")
			assert.Equal(t, tt.want, rr.Body.String(), "wrong response body")
		})
	}
}

func createReqListDB(r *req, t *testing.T) *http.Request {
	req, err := http.NewRequest("GET", r.url, nil)
	if err != nil {
		t.Fatal(err)
	}
	req.Header.Add("X-identityID", r.userHead)
	return req
}

func getListDBTestCase() []listDBCase {
	return []listDBCase{
		// TODO: Add test cases.
		{"one-day", &req{"http://localhost/recstorage/test/date?s=2018-11-08&e=2018-05-11", "user"},
			"{\"2018-05-11\":[\"2018-05-11/1/t.db\"]}", 200},
		{"one-day-more-db", &req{"http://localhost/recstorage/test/date?s=2018-05-08&e=2018-05-08", "user"},
			"{\"2018-05-08\":[\"2018-05-08/1/t.db\",\"2018-05-08/2/t.db\",\"2018-05-08/3/t.db\"]}", 200},
		{"two-day", &req{"http://localhost/recstorage/test/date?s=2018-05-08&e=2018-05-09", "user"},
			"{\"2018-05-08\":[\"2018-05-08/1/t.db\",\"2018-05-08/2/t.db\",\"2018-05-08/3/t.db\"],\"2018-05-09\":[\"2018-05-09/1/t.db\",\"2018-05-09/2/t.db\"]}", 200},
		{"more-day", &req{"http://localhost/recstorage/test/date?s=2018-05-06&e=2018-05-10", "user"},
			"{\"2018-05-06\":[\"2018-05-06/1/t.db\",\"2018-05-06/2/t.db\",\"2018-05-06/3/t.db\"],\"2018-05-08\":[\"2018-05-08/1/t.db\",\"2018-05-08/2/t.db\",\"2018-05-08/3/t.db\"],\"2018-05-09\":[\"2018-05-09/1/t.db\",\"2018-05-09/2/t.db\"],\"2018-05-10\":[\"2018-05-10/1/t.db\",\"2018-05-10/2/t.db\"]}", 200},
		{"more-day-invert", &req{"http://localhost/recstorage/test/date?s=2018-05-10&e=2018-05-06", "user"},
			"{\"2018-05-06\":[\"2018-05-06/1/t.db\",\"2018-05-06/2/t.db\",\"2018-05-06/3/t.db\"],\"2018-05-08\":[\"2018-05-08/1/t.db\",\"2018-05-08/2/t.db\",\"2018-05-08/3/t.db\"],\"2018-05-09\":[\"2018-05-09/1/t.db\",\"2018-05-09/2/t.db\"],\"2018-05-10\":[\"2018-05-10/1/t.db\",\"2018-05-10/2/t.db\"]}", 200},
		{"no-db-first-day", &req{"http://localhost/recstorage/test/date?s=2018-05-04&e=2018-05-10", "user"},
			"{\"2018-05-06\":[\"2018-05-06/1/t.db\",\"2018-05-06/2/t.db\",\"2018-05-06/3/t.db\"],\"2018-05-08\":[\"2018-05-08/1/t.db\",\"2018-05-08/2/t.db\",\"2018-05-08/3/t.db\"],\"2018-05-09\":[\"2018-05-09/1/t.db\",\"2018-05-09/2/t.db\"],\"2018-05-10\":[\"2018-05-10/1/t.db\",\"2018-05-10/2/t.db\"]}", 200},
		{"moer-day-no-db", &req{"http://localhost/recstorage/test/date?s=2018-05-01&e=2018-05-05", "user"},
			"{}", 200},
		{"no-db-first-last-day", &req{"http://localhost/recstorage/test/date?s=2018-05-05&e=2018-05-07", "user"},
			"{\"2018-05-06\":[\"2018-05-06/1/t.db\",\"2018-05-06/2/t.db\",\"2018-05-06/3/t.db\"]}", 200},
		{"miss-e", &req{"http://localhost/recstorage/test/date?s=2018-11-08", "user"},
			"parameter error\n", 400},
		{"miss-s", &req{"http://localhost/recstorage/test/date?e=2018-11-08", "user"},
			"parameter error\n", 400},
		{"miss-user", &req{"http://localhost/recstorage/test/date?s=2018-11-08&e=2018-11-08", ""},
			"401 permission denied\n", 401},
		{"miss-cameraid", &req{"http://localhost/recstorage/date?e=2018-11-08", "user"},
			"404 page not found\n", 404},
	}
}

func createMockList(ctrl *gomock.Controller) *mocks.MockLister {
	dates := []struct {
		prefix string   // input: filename
		list   []string // output: database list
	}{
		{"user/test/2018-05-06", []string{
			"user/test/2018-05-06/1/t.db",
			"user/test/2018-05-06/2/t.db",
			"user/test/2018-05-06/3/t.db"}},
		{"user/test/2018-05-08", []string{
			"user/test/2018-05-08/1/t.db",
			"user/test/2018-05-08/2/t.db",
			"user/test/2018-05-08/3/t.db"}},
		{"user/test/2018-05-09", []string{
			"user/test/2018-05-09/1/t.db",
			"user/test/2018-05-09/2/t.db"}},
		{"user/test/2018-05-10", []string{
			"user/test/2018-05-10/1/t.db",
			"user/test/2018-05-10/2/t.db"}},
		{"user/test/2018-05-11", []string{
			"user/test/2018-05-11/1/t.db"}},
	}
	mock := mocks.NewMockLister(ctrl)
	mock.EXPECT().List(gomock.Any()).DoAndReturn(func(f *storage.FileInfo) ([]string, error) {
		for _, dd := range dates {
			if f.FileName == dd.prefix {
				ret := make([]string, len(dd.list))
				copy(ret, dd.list)
				return ret, nil
			}
		}
		return nil, nil
	}).AnyTimes()
	return mock
}

type listDBCredCase struct {
	name       string
	cerd       *awsCred
	want       string
	wantStatus int
}

func TestListDBCrendential(t *testing.T) {
	queryReq := &req{"http://localhost/recstorage/test/date?s=2018-05-11&e=2018-05-11", "user"}
	// setup mock module
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	defaultService.cloud = &mockStorage{List: createMockCredList(mockCtrl)}
	router := mux.NewRouter().StrictSlash(true)
	router.Methods("GET").Path("/recstorage/{cameraid}/date").Name("testListDb").HandlerFunc(ListDb)
	tests := getListDBCredTestCase()
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			rr := httptest.NewRecorder()
			router.ServeHTTP(rr, addCredHeader(createReqListDB(queryReq, t), tt.cerd, t))
			assert.Equal(t, tt.wantStatus, rr.Code, "wrong response status")
			assert.Equal(t, tt.want, rr.Body.String(), "wrong response body")
		})
	}
}

func getListDBCredTestCase() []listDBCredCase {
	return []listDBCredCase{
		// TODO: Add test cases.
		{"ami-user-cred",
			&awsCred{"aaaaaaaaaaaaaaaaaaaaa",
				"bbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbb",
				""},
			"{\"2018-05-11\":[\"2018-05-11/1/t.db\"]}", 200},
		{"temp-cred",
			&awsCred{"tttttttttttttttttttt",
				"dddddddddddddddddddddddddddddddddddddddd",
				"yyyyyyyyyyyyyyyyyyyyyyyyyyyyyyyyyyyyyyyyyyyyyyyyyyyyyyyyyyyyyyyy"},
			"{\"2018-05-11\":[\"2018-05-11/1/t.db\"]}", 200},
		{"no-cred",
			&awsCred{"",
				"",
				""},
			"401 permission denied\n", 401},
		{"miss-id",
			&awsCred{"",
				"bbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbb",
				""},
			"401 permission denied\n", 401},
		{"miss-key",
			&awsCred{"tttttttttttttttttttt",
				"",
				""},
			"401 permission denied\n", 401},
	}
}

func createMockCredList(ctrl *gomock.Controller) *mocks.MockLister {
	dbList := []string{
		"user/test/2018-05-11/1/t.db",
	}
	mock := mocks.NewMockLister(ctrl)
	mock.EXPECT().List(gomock.Any()).DoAndReturn(func(f *storage.FileInfo) ([]string, error) {
		for _, dd := range mockCredentials {
			if f.AccessKeyID == dd.ID &&
				f.SecretKey == dd.Key &&
				f.SessionToken == dd.Token {
				ret := make([]string, len(dbList))
				copy(ret, dbList)
				return ret, nil
			}
		}
		return nil, fmt.Errorf("401 permission denied")
	}).AnyTimes()
	return mock
}

// utility funcs

func createReq(r *req, t *testing.T, method string) *http.Request {
	req, err := http.NewRequest(method, r.url, nil)
	if err != nil {
		t.Fatal(err)
	}
	return req
}

func createGetReq(r *req, t *testing.T) *http.Request {
	req := createReq(r, t, "GET")
	req.Header.Add("X-identityID", r.userHead)
	return req
}
func createPostReq(r *req, t *testing.T) *http.Request {
	req := createReq(r, t, "POST")
	req.Header.Add("X-identityID", r.userHead)
	return req
}
func addCredHeader(r *http.Request, c *awsCred, t *testing.T) *http.Request {
	r.Header.Add("X-accessKeyID", c.ID)
	r.Header.Add("X-secretKey", c.Key)
	r.Header.Add("X-sessionToken", c.Token)
	return r
}

func TestSetRegion(t *testing.T) {
	type args struct {
		r string
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		// TODO: Add test cases.
		{"change", args{"eu-central-1"}, "eu-central-1"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			SetRegion(tt.args.r)
			assert.Equal(t, tt.want, defaultService.region, "set region fail")
		})
	}
}

func TestServiceSetRegion(t *testing.T) {
	type args struct {
		r string
	}
	tests := []struct {
		name string
		s    *Service
		args args
		want *Service
	}{
		// TODO: Add test cases.
		{"new", &Service{}, args{"ap-southeast-2"}, &Service{region: "ap-southeast-2"}},
		{"change", &Service{region: "ap-southeast-1"}, args{"eu-central-1"}, &Service{region: "eu-central-1"}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.s.SetRegion(tt.args.r)
			assert.EqualValues(t, tt.want, got, "set region fail")
		})
	}
}

func TestSetBucket(t *testing.T) {
	type args struct {
		b string
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		// TODO: Add test cases.
		{"change", args{"aaaa-bbbb-cccccccccccccc-dddd"}, "aaaa-bbbb-cccccccccccccc-dddd"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			SetBucket(tt.args.b)
			assert.Equal(t, tt.want, defaultService.bucket, "set bucket fail")
		})
	}
}

func TestServiceSetBucket(t *testing.T) {
	type args struct {
		b string
	}
	tests := []struct {
		name string
		s    *Service
		args args
		want *Service
	}{
		// TODO: Add test cases.
		{"new", &Service{}, args{"aaaaa-bbb-ccccccc-dddd"}, &Service{bucket: "aaaaa-bbb-ccccccc-dddd"}},
		{"change", &Service{bucket: "xxxxx-hhh-ccccccc-yyyy"}, args{"aaaaa-eee-ccccccc-gggg"}, &Service{bucket: "aaaaa-eee-ccccccc-gggg"}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.s.SetBucket(tt.args.b)
			assert.Equal(t, tt.want, got, "set bucket fail")
		})
	}
}
