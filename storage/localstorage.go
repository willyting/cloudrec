package storage

import (
	"gachamachine/storage/iface"
	"io"
	"os"
	"path/filepath"
	"strings"
)

// SystemOperator ...
type SystemOperator interface {
	Open(*FileInfo) (io.ReadWriteCloser, error)
	OpenDir(*FileInfo) (iface.FolderOperator, error)
}

// LocalStroage ...
type LocalStroage struct {
}

// FileOperator ...
type FileOperator struct {
	SystemOperator
}

// LocalFileSystem ...
type LocalFileSystem struct {
}

// GetDownloader ...
func (s *LocalStroage) GetDownloader() Downloader {
	return &FileOperator{SystemOperator: &LocalFileSystem{}}
}

// GetUploader ...
func (s *LocalStroage) GetUploader() Uploader {
	return &FileOperator{SystemOperator: &LocalFileSystem{}}
}

// GetLister ...
func (s *LocalStroage) GetLister() Lister {
	return &FileOperator{SystemOperator: &LocalFileSystem{}}
}

// Download ...
func (s *FileOperator) Download(file *FileInfo, writeTo io.Writer) error {
	f, err := s.Open(file)
	if err != nil {
		return err
	}
	_, err = io.Copy(writeTo, f)
	f.Close()
	if err != nil {
		return err
	}
	return nil
}

// Upload ...
func (s *FileOperator) Upload(file *FileInfo, readFrom io.ReadCloser) error {
	f, err := s.Open(file)
	if err != nil {
		return err
	}
	_, err = io.Copy(f, readFrom)
	f.Close()
	readFrom.Close()
	if err != nil {
		return err
	}
	return nil
}

// List ...
func (s *FileOperator) List(file *FileInfo) ([]string, error) {
	f, err := s.OpenDir(file)
	if err != nil {
		return nil, err
	}
	files, err := f.Readdirnames(-1)
	if err != nil {
		return nil, err
	}
	filter := strings.Replace(file.FileName, "/", ".", -1)
	var list []string
	for _, ff := range files {
		if strings.Contains(ff, filter) && filepath.Ext(ff) == ".db" {
			list = append(list, strings.Replace(ff, ".", "/", 4))
		}
	}
	return list, nil
}

// Open return client on success, and return nil and an error on fail
func (s *LocalFileSystem) Open(file *FileInfo) (io.ReadWriteCloser, error) {
	return os.OpenFile("./storage/"+strings.Replace(file.FileName, "/", ".", -1), os.O_CREATE|os.O_RDWR, 0666)
}

// OpenDir return client on success, and return nil and an error on fail
func (s *LocalFileSystem) OpenDir(file *FileInfo) (iface.FolderOperator, error) {
	return os.Open("./storage/")
}
