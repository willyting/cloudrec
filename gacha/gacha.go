package gacha

import (
	"GaChaMachine/storage"
	"net/http"
)

var cloud storage.Storage

func setupStorage(upload storage.Uploader, dlwonload storage.Downloader) {
	cloud.Upload = upload
	cloud.Dlwonload = dlwonload
}

// GetRec ...
func GetRec(w http.ResponseWriter, r *http.Request) {
	cloud.Dlwonload.Download(&storage.FileInfo{}, w)
}
