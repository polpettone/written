package models

import "os"

type Document struct {
	Path string
	Info os.FileInfo
	Tags []string
}

func (document Document) String() string {
	return document.Info.Name()
}
