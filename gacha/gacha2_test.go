package gacha

//go:generate $GOPATH/bin/mockgen -destination src/gachamachine/mocks/mock_storage_uploader.go -package mocks gachamachine/storage Uploader
//go:generate $GOPATH/bin/mockgen -destination src/gachamachine/mocks/mock_storage_downloader.go -package mocks gachamachine/storage Downloader
//go:generate $GOPATH/bin/mockgen -destination src/gachamachine/mocks/mock_storage_lister.go -package mocks gachamachine/storage Lister
import (
	"bytes"
	"fmt"
	"gachamachine/mocks"
	"gachamachine/storage"
	"io"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/gorilla/mux"

	"github.com/golang/mock/gomock"
)

type mockStorage struct {
	Up   storage.Uploader
	Down storage.Downloader
	List storage.Lister
}

func (m *mockStorage) GetUploader() storage.Uploader {
	return m.Up
}
func (m *mockStorage) GetDownloader() storage.Downloader {
	return m.Down
}
func (m *mockStorage) GetLister() storage.Lister {
	return m.List
}

func TestGetRec_(t *testing.T) {
	router := mux.NewRouter().StrictSlash(true)
	router.Methods("GET").Path("/recstorage/{cameraid}").Name("test").HandlerFunc(GetRec)
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	mockStroage := mocks.NewMockDownloader(mockCtrl)
	testCloud := &mockStorage{Down: mockStroage}
	defaultService.cloud = testCloud
	mockStroage.EXPECT().Download(&storage.FileInfo{
		Region:   "ap-southeast-1",
		Bucket:   "ec4f9e12-5286-11e8-9c2d-fa7ae01bbebc",
		FileName: "user/test/test.txt",
	}, gomock.Any()).
		Return(nil).Do(func(f *storage.FileInfo, w io.Writer) {
		fmt.Fprint(w, "test")
	})
	req, err := http.NewRequest("GET", "http://localhost/recstorage/test?p=test.txt", nil)
	assert.NoError(t, err)
	req.Header.Add("X-identityID", "user")
	rr := httptest.NewRecorder()
	router.ServeHTTP(rr, req)
	assert.Equal(t, http.StatusOK, rr.Code, "response status incorrect")
	assert.Equal(t, "test", rr.Body.String(), "response body incorrect")
}

func TestPutRec_(t *testing.T) {
	// setup server handle
	router := mux.NewRouter().StrictSlash(true)
	router.Methods("POST").Path("/recstorage/{cameraid}").Name("test").HandlerFunc(PutRec)

	// setup mock storage
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	mockStroage := mocks.NewMockUploader(mockCtrl)
	testCloud := &mockStorage{Up: mockStroage}
	defaultService.cloud = testCloud
	mockStroage.EXPECT().Upload(&storage.FileInfo{
		Region:   "ap-southeast-1",
		Bucket:   "ec4f9e12-5286-11e8-9c2d-fa7ae01bbebc",
		FileName: "user/test/test.txt",
	}, gomock.Any()).
		Return(nil).Do(func(f *storage.FileInfo, r io.Reader) {
		// THEN: get ok response. storage will receive a new file on `S3:bucket/{userID}/{cameraID}/{filename}`
		stroageBuf, errF := ioutil.ReadAll(r)
		assert.NoError(t, errF)
		assert.Equal(t, "test", string(stroageBuf), "incorrect file contant")
	})

	// GIVE: cameraID = "test" and filename = "test.txt" in URL, userID="user" in header
	body := bytes.NewReader([]byte("test"))
	req, err := http.NewRequest("POST", "http://localhost/recstorage/test?p=test.txt", body)
	assert.NoError(t, err)
	req.Header.Add("X-identityID", "user")
	rr := httptest.NewRecorder()

	// WHEN: send a put file request
	router.ServeHTTP(rr, req)

	// THEN: get ok response. storage will receive a new file on `S3:bucket/{userID}/{cameraID}/{filename}`
	assert.Equal(t, http.StatusOK, rr.Code, "response status incorrect")
}

func TestGetRec_default(t *testing.T) {
	router := mux.NewRouter().StrictSlash(true)
	router.Methods("GET").Path("/recstorage/{cameraid}").Name("test").HandlerFunc(GetRec)
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	mockStroage := mocks.NewMockDownloader(mockCtrl)
	testCloud := &mockStorage{Down: mockStroage}
	defaultService.cloud = testCloud
	mockStroage.EXPECT().Download(&storage.FileInfo{
		Region:   "ap-southeast-1",
		Bucket:   "ec4f9e12-5286-11e8-9c2d-fa7ae01bbebc",
		FileName: "user/test/test.txt",
	}, gomock.Any()).Return(nil)
	req, err := http.NewRequest("GET", "http://localhost/recstorage/test?p=test.txt", nil)
	if err != nil {
		t.Fatal(err)
	}
	req.Header.Add("X-identityID", "user")
	rr := httptest.NewRecorder()
	router.ServeHTTP(rr, req)
	assert.Equal(t, http.StatusOK, rr.Code, "response status incorrect")
}

func TestGetRec_noUserID(t *testing.T) {
	router := mux.NewRouter().StrictSlash(true)
	router.Methods("GET").Path("/recstorage/{cameraid}").Name("test").HandlerFunc(GetRec)
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	mockStroage := mocks.NewMockDownloader(mockCtrl)
	testCloud := &mockStorage{Down: mockStroage}
	defaultService.cloud = testCloud
	mockStroage.EXPECT().Download(gomock.Any(), gomock.Any()).Return(nil).Times(0)
	req, err := http.NewRequest("GET", "http://localhost/recstorage/test?p=test.txt", nil)
	assert.NoError(t, err)
	rr := httptest.NewRecorder()
	router.ServeHTTP(rr, req)
	assert.Equal(t, http.StatusUnauthorized, rr.Code, "response status incorrect")
}
func TestPutRec_noUserID(t *testing.T) {
	router := mux.NewRouter().StrictSlash(true)
	router.Methods("POST").Path("/recstorage/{cameraid}").Name("test").HandlerFunc(PutRec)
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	mockStroage := mocks.NewMockUploader(mockCtrl)
	testCloud := &mockStorage{Up: mockStroage}
	defaultService.cloud = testCloud
	mockStroage.EXPECT().Upload(gomock.Any(), gomock.Any()).Return(nil).Times(0)
	req, err := http.NewRequest("POST", "http://localhost/recstorage/test?p=test.txt", nil)
	assert.NoError(t, err)
	rr := httptest.NewRecorder()
	router.ServeHTTP(rr, req)
	assert.Equal(t, http.StatusUnauthorized, rr.Code, "response status incorrect")
}
