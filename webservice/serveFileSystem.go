package webservice

import (
	"embed"
	"fmt"
	"net/http"
	"path"
)

//go:embed www/regate/dist
var systemFileInternal embed.FS

// Create Service file for static page css+images+js
type ServeFileSystem struct {
	http.FileSystem

	prefix string
}

func (sfs *ServeFileSystem) Exists(prefix string, pathFile string) bool {
	p := path.Join(sfs.prefix, pathFile)

	fmt.Println("Check: ", p)
	_, err := sfs.FileSystem.Open(p)
	if err != nil {
		return false
	}
	return true
}

func (sfs *ServeFileSystem) Open(pathFile string) (http.File, error) {
	return sfs.FileSystem.Open(path.Join(sfs.prefix, pathFile))
}

// ServerFile into memory
func NewServeFileSystem(fs embed.FS, prefix string) *ServeFileSystem {
	var sfs ServeFileSystem
	sfs.prefix = prefix
	sfs.FileSystem = http.FS(fs)
	return &sfs
}
