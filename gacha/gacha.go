package gacha

import (
	"GaChaMachine/storage"
	"net/http"

	"github.com/gorilla/mux"
)

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
