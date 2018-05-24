package storage

import "io"

// FileInfo ...
type FileInfo struct {
	Region       string
	Bucket       string
	FileName     string
	AccessKeyID  string
	SecretKey    string
	SessionToken string
}

// Downloader ...
type Downloader interface {
	Download(*FileInfo, io.Writer) error
}

// Uploader ...
type Uploader interface {
	Upload(*FileInfo, io.ReadCloser) error
}

// Lister ...
type Lister interface {
	List(*FileInfo) ([]string, error)
}

// Storage ...
type Storage interface {
	GetUploader() Uploader
	GetDownloader() Downloader
	GetLister() Lister
}
