package gacha

import (
	"GaChaMachine/storage"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/gorilla/mux"
)

//go:generate $GOPATH/bin/mockgen -destination src/GaChaMachine/mocks/mock_s3_client.go -package mocks github.com/aws/aws-sdk-go/service/s3/s3iface S3API

var cloud storage.Storage

// GetRec ...
func GetRec(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	cameraID := vars["cameraid"]
	userID := r.Header.Get("X-identityID")
	filePath := r.URL.Query().Get("p")
	cloud.GetDownloader().Download(&storage.FileInfo{
		FileName: userID + "/" + cameraID + "/" + filePath,
	}, w)
}

// PutRec ...
func PutRec(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	cameraID := vars["cameraid"]
	userID := r.Header.Get("X-identityID")
	filePath := r.URL.Query().Get("p")
	cloud.GetUploader().Upload(&storage.FileInfo{
		FileName: userID + "/" + cameraID + "/" + filePath,
	}, r.Body)
	r.Body.Close()
}

// ListDb ...
func ListDb(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	cameraID := vars["cameraid"]
	userID := r.Header.Get("X-identityID")
	startDate := r.URL.Query().Get("s")
	endDate := r.URL.Query().Get("e")
	lister := cloud.GetLister()

	start, err := time.Parse("2006-01-02", startDate)
	if err != nil {
		fmt.Fprint(w, err)
		return
	}
	end, err := time.Parse("2006-01-02", endDate)
	if err != nil {
		fmt.Fprint(w, err)
		return
	}
	if start.After(end) {
		tmp := start
		start = end
		end = tmp
	}
	recList := make(map[string][]string)
	base := userID + "/" + cameraID + "/"
	for !start.After(end) {
		curDate := start.Format("2006-01-02")
		list, err := lister.List(&storage.FileInfo{
			FileName: base + curDate,
		})
		if err != nil {
			fmt.Fprint(w, err)
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
