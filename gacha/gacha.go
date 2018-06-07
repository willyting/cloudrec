package gacha

import (
	"encoding/json"
	"fmt"
	"gachamachine/machine"
	"gachamachine/storage"
	"net/http"
	"time"

	"github.com/gorilla/mux"
)

// Service ...
type Service struct {
	cloud  storage.Storage
	region string
	bucket string
}

var defaultService = &Service{
	cloud:  &storage.S3Stroage{},
	region: "ap-southeast-1",
	bucket: "ec4f9e12-5286-11e8-9c2d-fa7ae01bbebc",
}

// GetHandlers ...
func GetHandlers() []machine.Route {
	return defaultService.GetHandlers()
}

// GetHandlers ...
func (s *Service) GetHandlers() []machine.Route {
	return []machine.Route{
		{Name: "playback", Method: "GET", Pattern: "/recstorage/{cameraid}", HandlerFunc: s.GetRec},
		{Name: "record-post", Method: "POST", Pattern: "/recstorage/{cameraid}", HandlerFunc: s.PutRec},
		{Name: "record-put", Method: "PUT", Pattern: "/recstorage/{cameraid}", HandlerFunc: s.PutRec},
		{Name: "date-query", Method: "GET", Pattern: "/recstorage/{cameraid}/date", HandlerFunc: s.ListDb},
	}
}

// SetRegion ...
func SetRegion(r string) {
	defaultService.SetRegion(r)
}

// SetRegion ...
func (s *Service) SetRegion(r string) *Service {
	s.region = r
	return s
}

// SetBucket ...
func SetBucket(b string) {
	defaultService.SetBucket(b)
}

// SetBucket ...
func (s *Service) SetBucket(b string) *Service {
	s.bucket = b
	return s
}

// SetStorage ...
func SetStorage(c storage.Storage) {
	defaultService.SetStorage(c)
}

// SetStorage ...
func (s *Service) SetStorage(c storage.Storage) *Service {
	s.cloud = c
	return s
}

// GetRec ...
func GetRec(w http.ResponseWriter, r *http.Request) {
	defaultService.GetRec(w, r)
}

// GetRec ...
func (s *Service) GetRec(w http.ResponseWriter, r *http.Request) {
	base, err := getUserCamPath(w, r)
	if err != nil {
		return
	}
	authInfo := s.infoParser(r)
	authInfo.FileName = base + r.URL.Query().Get("p")
	err = s.cloud.GetDownloader().Download(authInfo, w)
	if err != nil {
		http.Error(w, "401 permission denied", 401)
		return
	}
}

// PutRec ...
func PutRec(w http.ResponseWriter, r *http.Request) {
	defaultService.PutRec(w, r)
}

// PutRec ...
func (s *Service) PutRec(w http.ResponseWriter, r *http.Request) {
	base, err := getUserCamPath(w, r)
	if err != nil {
		return
	}
	authInfo := s.infoParser(r)
	authInfo.FileName = base + r.URL.Query().Get("p")
	err = s.cloud.GetUploader().Upload(authInfo, r.Body)
	if err != nil {
		http.Error(w, "401 permission denied", 401)
		return
	}
}

// ListDb ...
func ListDb(w http.ResponseWriter, r *http.Request) {
	defaultService.ListDb(w, r)
}

// ListDb ...
func (s *Service) ListDb(w http.ResponseWriter, r *http.Request) {
	base, err := getUserCamPath(w, r)
	if err != nil {
		return
	}
	start, err := time.Parse("2006-01-02", r.URL.Query().Get("s"))
	if err != nil {
		http.Error(w, "parameter error", 400)
		return
	}
	end, err := time.Parse("2006-01-02", r.URL.Query().Get("e"))
	if err != nil {
		http.Error(w, "parameter error", 400)
		return
	}
	if start.After(end) {
		tmp := start
		start = end
		end = tmp
	}
	authInfo := s.infoParser(r)
	lister := s.cloud.GetLister()
	recList := make(map[string][]string)
	for !start.After(end) {
		curDate := start.Format("2006-01-02")
		authInfo.FileName = base + curDate
		list, err := lister.List(authInfo)
		if err != nil {
			http.Error(w, "401 permission denied", 401)
			return
		}
		if list != nil {
			for i := range list {
				list[i] = list[i][len(base):len(list[i])]
			}
			recList[curDate] = list
		}
		start = start.AddDate(0, 0, 1)
	}
	result, err := json.Marshal(recList)
	fmt.Fprint(w, string(result))
}

func (s *Service) infoParser(r *http.Request) *storage.FileInfo {
	f := s.newDefaultInfo()
	setCredentialInfo(f, r)
	return f
}

func (s *Service) newDefaultInfo() *storage.FileInfo {
	return &storage.FileInfo{
		Region: s.region,
		Bucket: s.bucket,
	}
}

func getUserCamPath(w http.ResponseWriter, r *http.Request) (string, error) {
	vars := mux.Vars(r)
	cameraID := vars["cameraid"]
	userID := r.Header.Get("X-identityID")
	if len(userID) == 0 {
		http.Error(w, "401 permission denied", 401)
		return "", fmt.Errorf("401 permission denied")
	}
	return userID + "/" + cameraID + "/", nil
}

func setCredentialInfo(f *storage.FileInfo, r *http.Request) {
	f.AccessKeyID = r.Header.Get("X-accessKeyID")
	f.SecretKey = r.Header.Get("X-secretKey")
	f.SessionToken = r.Header.Get("X-sessionToken")
}
