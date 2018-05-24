package gacha

//go:generate $GOPATH/bin/mockgen -destination /src/GaChaMachine/mocks/mock_storage_uploader.go -package mocks GaChaMachine/storage Uploader
//go:generate $GOPATH/bin/mockgen -destination /src/GaChaMachine/mocks/mock_storage_downloader.go -package mocks GaChaMachine/storage Downloader
import (
	"GaChaMachine/mocks"
	"GaChaMachine/storage"
	"bytes"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"

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
	cloud = testCloud
	mockStroage.EXPECT().Download(&storage.FileInfo{
		FileName: "user/test/test.txt",
	}, gomock.Any()).
		Return(nil).Do(func(f *storage.FileInfo, w io.Writer) {
		fmt.Fprint(w, "test")
	})
	req, err := http.NewRequest("GET", "http://localhost/recstorage/test?p=test.txt", nil)
	if err != nil {
		t.Fatal(err)
	}
	req.Header.Add("X-identityID", "user")
	rr := httptest.NewRecorder()
	router.ServeHTTP(rr, req)
	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %d want %d",
			status, http.StatusOK)
	}
	expected := "test"
	buf, err := ioutil.ReadAll(rr.Body)
	if string(buf) != expected {
		t.Errorf("handler returned unexpected body: got %s want %s",
			string(buf), expected)
	}
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
	cloud = testCloud
	mockStroage.EXPECT().Upload(&storage.FileInfo{
		FileName: "user/test/test.txt",
	}, gomock.Any()).
		Return(nil).Do(func(f *storage.FileInfo, r io.Reader) {
		// THEN: get ok response. storage will receive a new file on `S3:bucket/{userID}/{cameraID}/{filename}`
		expected := "test"
		stroageBuf, errF := ioutil.ReadAll(r)
		if errF != nil {
			t.Errorf(errF.Error())
		}
		if string(stroageBuf) != expected {
			t.Errorf("incorrect file got: %s, want: %s", string(stroageBuf), expected)
		}
	})

	// GIVE: cameraID = "test" and filename = "test.txt" in URL, userID="user" in header
	body := bytes.NewReader([]byte("test"))
	req, err := http.NewRequest("POST", "http://localhost/recstorage/test?p=test.txt", body)
	if err != nil {
		t.Fatal(err)
	}
	req.Header.Add("X-identityID", "user")
	rr := httptest.NewRecorder()

	// WHEN: send a put file request
	router.ServeHTTP(rr, req)

	// THEN: get ok response. storage will receive a new file on `S3:bucket/{userID}/{cameraID}/{filename}`
	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %d want %d",
			status, http.StatusOK)
	}
}
