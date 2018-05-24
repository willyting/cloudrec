package gacha

import (
	"GaChaMachine/machine"
	"GaChaMachine/storage"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/gorilla/mux"
)

var cloud storage.Storage

// GetHandlers ...
func GetHandlers() []machine.Route {
	return []machine.Route{
		{"playback", "GET", "/recstorage/{cameraid}", GetRec},
		{"record-post", "POST", "/recstorage/{cameraid}", PutRec},
		{"record-put", "PUT", "/recstorage/{cameraid}", PutRec},
		{"date-query", "GET", "/recstorage/{cameraid}/date", ListDb},
	}
}

// GetRec ...
func GetRec(w http.ResponseWriter, r *http.Request) {
	base, err := getUserCamPath(w, r)
	if err != nil {
		return
	}
	filePath := r.URL.Query().Get("p")
	cloud.GetDownloader().Download(&storage.FileInfo{
		FileName: base + filePath,
	}, w)
}

// PutRec ...
func PutRec(w http.ResponseWriter, r *http.Request) {
	base, err := getUserCamPath(w, r)
	if err != nil {
		return
	}
	filePath := r.URL.Query().Get("p")
	cloud.GetUploader().Upload(&storage.FileInfo{
		FileName: base + filePath,
	}, r.Body)
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
	lister := cloud.GetLister()
	recList := make(map[string][]string)
	for !start.After(end) {
		curDate := start.Format("2006-01-02")
		list, _ := lister.List(&storage.FileInfo{
			FileName: base + curDate,
		})
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
