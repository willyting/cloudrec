package gacha

import (
	"GaChaMachine/storage"
	"fmt"
	"net/http"

	"github.com/gorilla/mux"
)

//go:generate $GOPATH/bin/mockgen -destination src/GaChaMachine/mocks/mock_s3_client.go -package mocks github.com/aws/aws-sdk-go/service/s3/s3iface S3API

var cloud storage.Storage

// func setupStorage(upload storage.Uploader, dlwonload storage.Downloader) {
// 	cloud.Upload = upload
// 	cloud.Dlwonload = dlwonload
// }

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
	fmt.Fprint(w, "tbd")
}
