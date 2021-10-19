package models

import "os"

type Document struct {
	Path string
	Info os.FileInfo
}

func (document Document) String() string {
	return document.Info.Name()
}
