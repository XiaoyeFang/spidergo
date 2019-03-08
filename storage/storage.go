package storage

import "io"

type Storage interface {
	CopyFile(localPath, remotePath string) error
	CopyFileByUrl(srcLink, remotePath string) error
	CopyFileByReader(r io.Reader, fileType, size, remotePath string) (string, error)
	LsDir(path string) ([]string, error)
	CheckExist(key string) bool
	RemoveDir(path string) error
}