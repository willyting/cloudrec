package gacha

import (
	"GaChaMachine/mocks"
	"GaChaMachine/storage"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/gorilla/mux"
)

type req struct {
	url      string
	userHead string
}

func TestGetRec(t *testing.T) {
	type args struct {
		w http.ResponseWriter
		r *http.Request
	}
	tests := []struct {
		name string
		args args
	}{
	// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			GetRec(tt.args.w, tt.args.r)
		})
	}
}

func TestPutRec(t *testing.T) {
	type args struct {
		w http.ResponseWriter
		r *http.Request
	}
	tests := []struct {
		name string
		args args
	}{
	// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			PutRec(tt.args.w, tt.args.r)
		})
	}
}

func TestListDb(t *testing.T) {
	// setup mock module
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	cloud = &mockStorage{List: createMockList(mockCtrl)}
	router := mux.NewRouter().StrictSlash(true)
	router.Methods("GET").Path("/recstorage/{cameraid}/date").Name("testListDb").HandlerFunc(ListDb)

	tests := []struct {
		name       string
		req        req // input: info to create http req
		want       string
		wantStatus int
	}{
		// TODO: Add test cases.
		{"on-day", req{"http://localhost/recstorage/test/date?s=2018-05-08&e=2018-05-08", "user"},
			"{\"2018-05-08\":[\"2018-05-08/001/2018-05-08.db\",\"2018-05-08/002/2018-05-08.db\",\"2018-05-08/003/2018-05-08.db\"]}", 200},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			rr := httptest.NewRecorder()
			router.ServeHTTP(rr, createReq(&tt.req, t))
			if status := rr.Code; status != tt.wantStatus {
				t.Errorf("handler returned wrong status code: got %d want %d", status, tt.wantStatus)
			}
			if gotBody := rr.Body.String(); gotBody != tt.want {
				t.Errorf("handler returned unexpected body: got %s want %s", gotBody, tt.want)
			}
		})
	}
}

func createReq(r *req, t *testing.T) *http.Request {
	req, err := http.NewRequest("GET", r.url, nil)
	if err != nil {
		t.Fatal(err)
	}
	req.Header.Add("X-identityID", r.userHead)
	return req
}

func createMockList(ctrl *gomock.Controller) *mocks.MockLister {
	mock := mocks.NewMockLister(ctrl)
	dates := []struct {
		prefix string   // input: filename
		list   []string // output: database list
	}{
		{"user/test/2018-05-06", []string{
			"user/test/2018-05-06/001/2018-05-06.db",
			"user/test/2018-05-06/002/2018-05-06.db",
			"user/test/2018-05-06/003/2018-05-06.db"}},
		{"user/test/2018-05-08", []string{
			"user/test/2018-05-08/001/2018-05-08.db",
			"user/test/2018-05-08/002/2018-05-08.db",
			"user/test/2018-05-08/003/2018-05-08.db"}},
		{"user/test/2018-05-09", []string{
			"user/test/2018-05-09/001/2018-05-09.db",
			"user/test/2018-05-09/002/2018-05-09.db"}},
		{"user/test/2018-05-10", []string{
			"user/test/2018-05-10/001/2018-05-10.db",
			"user/test/2018-05-10/002/2018-05-10.db"}},
		{"user/test/2018-05-11", []string{
			"user/test/2018-05-11/002/2018-05-11.db"}},
	}
	for _, dd := range dates {
		mock.EXPECT().List(&storage.FileInfo{
			FileName: dd.prefix,
		}).Return(dd.list, nil).AnyTimes()
	}
	mock.EXPECT().List(gomock.Any()).Return(nil, nil).AnyTimes()
	return mock
}
