package gacha

import (
	"GaChaMachine/mocks"
	"GaChaMachine/storage"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gorilla/mux"

	"github.com/golang/mock/gomock"
)

func TestGetRec(t *testing.T) {
	router := mux.NewRouter().StrictSlash(true)
	router.Methods("GET").Path("/recstorage/{cameraid}").Name("test").HandlerFunc(GetRec)
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	mockUploader := mocks.NewMockDownloader(mockCtrl)
	cloud.Dlwonload = mockUploader
	mockUploader.EXPECT().Download(&storage.FileInfo{
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
