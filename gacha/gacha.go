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

var cloud storage.Storage

// GetHandlers ...
func GetHandlers() []machine.Route {
	return []machine.Route{
		{Name: "playback", Method: "GET", Pattern: "/recstorage/{cameraid}", HandlerFunc: GetRec},
		{Name: "record-post", Method: "POST", Pattern: "/recstorage/{cameraid}", HandlerFunc: PutRec},
		{Name: "record-put", Method: "PUT", Pattern: "/recstorage/{cameraid}", HandlerFunc: PutRec},
		{Name: "date-query", Method: "GET", Pattern: "/recstorage/{cameraid}/date", HandlerFunc: ListDb},
	}
}

// GetRec ...
func GetRec(w http.ResponseWriter, r *http.Request) {
	base, err := getUserCamPath(w, r)
	if err != nil {
		return
	}
	authInfo := getCredentialInfo(r)
	authInfo.FileName = base + r.URL.Query().Get("p")
	err = cloud.GetDownloader().Download(authInfo, w)
	if err != nil {
		http.Error(w, "401 permission denied", 401)
		return
	}
}

// PutRec ...
func PutRec(w http.ResponseWriter, r *http.Request) {
	base, err := getUserCamPath(w, r)
	if err != nil {
		return
	}
	authInfo := getCredentialInfo(r)
	authInfo.FileName = base + r.URL.Query().Get("p")
	err = cloud.GetUploader().Upload(authInfo, r.Body)
	if err != nil {
		http.Error(w, "401 permission denied", 401)
		return
	}
}

// ListDb ...
func ListDb(w http.ResponseWriter, r *http.Request) {
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
	authInfo := getCredentialInfo(r)
	lister := cloud.GetLister()
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

func getCredentialInfo(r *http.Request) *storage.FileInfo {
	return &storage.FileInfo{
		AccessKeyID:  r.Header.Get("X-accessKeyID"),
		SecretKey:    r.Header.Get("X-secretKey"),
		SessionToken: r.Header.Get("X-sessionToken"),
	}
}
