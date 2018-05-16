package gacha

import (
	"GaChaMachine/mocks"
	"GaChaMachine/storage"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/golang/mock/gomock"
)

func TestGetRec(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	mockUploader := mocks.NewMockDownloader(mockCtrl)
	cloud.Dlwonload = mockUploader
	mockUploader.EXPECT().Download(gomock.Any(), gomock.Any()).
		Return(nil).Do(func(f *storage.FileInfo, w io.Writer) {
		fmt.Fprint(w, "test123")
	})
	req, err := http.NewRequest("GET", "/recstorage/test?p=test", nil)
	if err != nil {
		t.Fatal(err)
	}
	req.Header.Add("X-identityID", "user")
	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(GetRec)
	handler.ServeHTTP(rr, req)
	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusOK)
	}
	expected := "test"
	if rr.Body.String() != expected {
		t.Errorf("handler returned unexpected body: got %v want %v",
			rr.Body.String(), expected)
	}
}
