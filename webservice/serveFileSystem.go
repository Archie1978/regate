package webservice

import (
	"embed"
	"net/http"
	"path"
)

/*

Service web file



*/

//go:embed www/regate/dist
var systemFileInternal embed.FS

// Create Service file for static page css+images+js
type ServeFileSystem struct {
	http.FileSystem

	prefix string
}

// check Exists file into  memory
func (sfs *ServeFileSystem) Exists(prefix string, pathFile string) bool {
	p := path.Join(sfs.prefix, pathFile)

	_, err := sfs.FileSystem.Open(p)
	if err != nil {
		return false
	}
	return true
}

// Open file into mémory
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
