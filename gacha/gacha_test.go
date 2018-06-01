package gacha

//go:generate $GOPATH/bin/mockgen -destination /src/gachamachine/mocks/mock_storage_uploader.go -package mocks gachamachine/storage Uploader
//go:generate $GOPATH/bin/mockgen -destination /src/gachamachine/mocks/mock_storage_downloader.go -package mocks gachamachine/storage Downloader
//go:generate $GOPATH/bin/mockgen -destination src/gachamachine/mocks/mock_storage_lister.go -package mocks gachamachine/storage Lister
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
type aws_Cred struct {
	Id    string // input: accessKeyID
	Key   string // input: secretKey
	Token string // input: filename
}

var mockCredentials = []aws_Cred{
	{"aaaaaaaaaaaaaaaaaaaaa",
		"bbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbb",
		""},
	{"tttttttttttttttttttt",
		"dddddddddddddddddddddddddddddddddddddddd",
		"yyyyyyyyyyyyyyyyyyyyyyyyyyyyyyyyyyyyyyyyyyyyyyyyyyyyyyyyyyyyyyyy"},
}

type rec_cred_Case struct {
	name       string
	cerd       *aws_Cred
	want       string
	wantStatus int
}

func TestGetRec_Crendential(t *testing.T) {
	queryReq := &req{"http://localhost/recstorage/test?p=test.txt", "user"}
	// setup mock module
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	cloud = &mockStorage{Down: createMockCredGetRec(mockCtrl)}
	router := mux.NewRouter().StrictSlash(true)
	router.Methods("GET").Path("/recstorage/{cameraid}").Name("test").HandlerFunc(GetRec)
	tests := getGetRec_Cred_TestCase()
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			rr := httptest.NewRecorder()
			router.ServeHTTP(rr, createReq_GetRec(queryReq, tt.cerd, t))
			assert.Equal(t, tt.wantStatus, rr.Code, "wrong response status")
			assert.Equal(t, tt.want, rr.Body.String(), "wrong response body")
		})
	}
}

func createReq_GetRec(r *req, c *aws_Cred, t *testing.T) *http.Request {
	return addCredHeader(createGetReq(r, t), c, t)
}

func getGetRec_Cred_TestCase() []rec_cred_Case {
	return []rec_cred_Case{
		// TODO: Add test cases.
		{"ami-user-cred",
			&aws_Cred{"aaaaaaaaaaaaaaaaaaaaa",
				"bbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbb",
				""},
			"test", 200},
		{"temp-cred",
			&aws_Cred{"tttttttttttttttttttt",
				"dddddddddddddddddddddddddddddddddddddddd",
				"yyyyyyyyyyyyyyyyyyyyyyyyyyyyyyyyyyyyyyyyyyyyyyyyyyyyyyyyyyyyyyyy"},
			"test", 200},
		{"fail-cred",
			&aws_Cred{"aaaaaaaaaaaaaaaaaa",
				"bbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbb",
				""},
			"401 permission denied\n", 401},
		{"no-cred",
			&aws_Cred{"",
				"",
				""},
			"401 permission denied\n", 401},
		{"miss-id",
			&aws_Cred{"",
				"bbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbb",
				""},
			"401 permission denied\n", 401},
		{"miss-key",
			&aws_Cred{"tttttttttttttttttttt",
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
				if f.AccessKeyID == dd.Id &&
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
	cloud = &mockStorage{Up: createMockCredPutRec(mockCtrl)}
	router := mux.NewRouter().StrictSlash(true)
	router.Methods("POST").Path("/recstorage/{cameraid}").Name("test").HandlerFunc(PutRec)
	tests := getGetRec_Cred_TestCase()
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			rr := httptest.NewRecorder()
			router.ServeHTTP(rr, createReq_PutRec(queryReq, tt.cerd, t))
			assert.Equal(t, tt.wantStatus, rr.Code, "wrong response status")
		})
	}
}
func createMockCredPutRec(ctrl *gomock.Controller) *mocks.MockUploader {
	mock := mocks.NewMockUploader(ctrl)
	mock.EXPECT().Upload(gomock.Any(), gomock.Any()).
		DoAndReturn(func(f *storage.FileInfo, r io.Reader) error {
			for _, dd := range mockCredentials {
				if f.AccessKeyID == dd.Id &&
					f.SecretKey == dd.Key &&
					f.SessionToken == dd.Token {
					return nil
				}
			}
			return fmt.Errorf("401 permission denied")
		}).AnyTimes()
	return mock
}
func createReq_PutRec(r *req, c *aws_Cred, t *testing.T) *http.Request {
	return addCredHeader(createPostReq(r, t), c, t)
}

type listDB_Case struct {
	name       string
	req        *req // input: info to create http req
	want       string
	wantStatus int
}

func TestListDb(t *testing.T) {
	// setup mock module
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	cloud = &mockStorage{List: createMockList(mockCtrl)}
	router := mux.NewRouter().StrictSlash(true)
	router.Methods("GET").Path("/recstorage/{cameraid}/date").Name("testListDb").HandlerFunc(ListDb)
	tests := getListDB_TestCase()
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			rr := httptest.NewRecorder()
			router.ServeHTTP(rr, createReq_ListDb(tt.req, t))
			assert.Equal(t, tt.wantStatus, rr.Code, "wrong response status")
			assert.Equal(t, tt.want, rr.Body.String(), "wrong response body")
		})
	}
}

func createReq_ListDb(r *req, t *testing.T) *http.Request {
	req, err := http.NewRequest("GET", r.url, nil)
	if err != nil {
		t.Fatal(err)
	}
	req.Header.Add("X-identityID", r.userHead)
	return req
}

func getListDB_TestCase() []listDB_Case {
	return []listDB_Case{
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

type listDB_Cred_Case struct {
	name       string
	cerd       *aws_Cred
	want       string
	wantStatus int
}

func TestListDb_Crendential(t *testing.T) {
	queryReq := &req{"http://localhost/recstorage/test/date?s=2018-05-11&e=2018-05-11", "user"}
	// setup mock module
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	cloud = &mockStorage{List: createMockCredList(mockCtrl)}
	router := mux.NewRouter().StrictSlash(true)
	router.Methods("GET").Path("/recstorage/{cameraid}/date").Name("testListDb").HandlerFunc(ListDb)
	tests := getListDB_Cred_TestCase()
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			rr := httptest.NewRecorder()
			router.ServeHTTP(rr, addCredHeader(createReq_ListDb(queryReq, t), tt.cerd, t))
			assert.Equal(t, tt.wantStatus, rr.Code, "wrong response status")
			assert.Equal(t, tt.want, rr.Body.String(), "wrong response body")
		})
	}
}

func getListDB_Cred_TestCase() []listDB_Cred_Case {
	return []listDB_Cred_Case{
		// TODO: Add test cases.
		{"ami-user-cred",
			&aws_Cred{"aaaaaaaaaaaaaaaaaaaaa",
				"bbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbb",
				""},
			"{\"2018-05-11\":[\"2018-05-11/1/t.db\"]}", 200},
		{"temp-cred",
			&aws_Cred{"tttttttttttttttttttt",
				"dddddddddddddddddddddddddddddddddddddddd",
				"yyyyyyyyyyyyyyyyyyyyyyyyyyyyyyyyyyyyyyyyyyyyyyyyyyyyyyyyyyyyyyyy"},
			"{\"2018-05-11\":[\"2018-05-11/1/t.db\"]}", 200},
		{"no-cred",
			&aws_Cred{"",
				"",
				""},
			"401 permission denied\n", 401},
		{"miss-id",
			&aws_Cred{"",
				"bbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbb",
				""},
			"401 permission denied\n", 401},
		{"miss-key",
			&aws_Cred{"tttttttttttttttttttt",
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
			if f.AccessKeyID == dd.Id &&
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
func addCredHeader(r *http.Request, c *aws_Cred, t *testing.T) *http.Request {
	r.Header.Add("X-accessKeyID", c.Id)
	r.Header.Add("X-secretKey", c.Key)
	r.Header.Add("X-sessionToken", c.Token)
	return r
}
