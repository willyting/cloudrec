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
		{"one-day", req{"http://localhost/recstorage/test/date?s=2018-11-08&e=2018-05-11", "user"},
			"{\"2018-05-11\":[\"2018-05-11/1/t.db\"]}", 200},
		{"one-day-more-db", req{"http://localhost/recstorage/test/date?s=2018-05-08&e=2018-05-08", "user"},
			"{\"2018-05-08\":[\"2018-05-08/1/t.db\",\"2018-05-08/2/t.db\",\"2018-05-08/3/t.db\"]}", 200},
		{"two-day", req{"http://localhost/recstorage/test/date?s=2018-05-08&e=2018-05-09", "user"},
			"{\"2018-05-08\":[\"2018-05-08/1/t.db\",\"2018-05-08/2/t.db\",\"2018-05-08/3/t.db\"],\"2018-05-09\":[\"2018-05-09/1/t.db\",\"2018-05-09/2/t.db\"]}", 200},
		{"more-day", req{"http://localhost/recstorage/test/date?s=2018-05-06&e=2018-05-10", "user"},
			"{\"2018-05-06\":[\"2018-05-06/1/t.db\",\"2018-05-06/2/t.db\",\"2018-05-06/3/t.db\"],\"2018-05-08\":[\"2018-05-08/1/t.db\",\"2018-05-08/2/t.db\",\"2018-05-08/3/t.db\"],\"2018-05-09\":[\"2018-05-09/1/t.db\",\"2018-05-09/2/t.db\"],\"2018-05-10\":[\"2018-05-10/1/t.db\",\"2018-05-10/2/t.db\"]}", 200},
		{"more-day-invert", req{"http://localhost/recstorage/test/date?s=2018-05-10&e=2018-05-06", "user"},
			"{\"2018-05-06\":[\"2018-05-06/1/t.db\",\"2018-05-06/2/t.db\",\"2018-05-06/3/t.db\"],\"2018-05-08\":[\"2018-05-08/1/t.db\",\"2018-05-08/2/t.db\",\"2018-05-08/3/t.db\"],\"2018-05-09\":[\"2018-05-09/1/t.db\",\"2018-05-09/2/t.db\"],\"2018-05-10\":[\"2018-05-10/1/t.db\",\"2018-05-10/2/t.db\"]}", 200},
		{"no-db-first-day", req{"http://localhost/recstorage/test/date?s=2018-05-04&e=2018-05-10", "user"},
			"{\"2018-05-06\":[\"2018-05-06/1/t.db\",\"2018-05-06/2/t.db\",\"2018-05-06/3/t.db\"],\"2018-05-08\":[\"2018-05-08/1/t.db\",\"2018-05-08/2/t.db\",\"2018-05-08/3/t.db\"],\"2018-05-09\":[\"2018-05-09/1/t.db\",\"2018-05-09/2/t.db\"],\"2018-05-10\":[\"2018-05-10/1/t.db\",\"2018-05-10/2/t.db\"]}", 200},
		{"moer-day-no-db", req{"http://localhost/recstorage/test/date?s=2018-05-01&e=2018-05-05", "user"},
			"{}", 200},
		{"no-db-first-last-day", req{"http://localhost/recstorage/test/date?s=2018-05-05&e=2018-05-07", "user"},
			"{\"2018-05-06\":[\"2018-05-06/1/t.db\",\"2018-05-06/2/t.db\",\"2018-05-06/3/t.db\"]}", 200},
		{"miss-e", req{"http://localhost/recstorage/test/date?s=2018-11-08", "user"},
			"parameter error\n", 400},
		{"miss-s", req{"http://localhost/recstorage/test/date?e=2018-11-08", "user"},
			"parameter error\n", 400},
		{"miss-user", req{"http://localhost/recstorage/test/date?s=2018-11-08&e=2018-11-08", ""},
			"401 permission denied\n", 401},
		{"miss-cameraid", req{"http://localhost/recstorage/date?e=2018-11-08", "user"},
			"404 page not found\n", 404},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			rr := httptest.NewRecorder()
			router.ServeHTTP(rr, createReq(&tt.req, t))
			if status := rr.Code; status != tt.wantStatus {
				t.Errorf("handler returned wrong status code: got %d want %d", status, tt.wantStatus)
			}
			if gotBody := rr.Body.String(); gotBody != tt.want {
				t.Errorf("handler returned unexpected body: got \n%s \nwant \n%s", gotBody, tt.want)
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
