package models

import (
	"fmt"
	"os"
)

type Document struct {
	Path string
	Info os.FileInfo
	Tags []string
	Content string
}

func (document Document) String() string {
	name := ""
	if document.Info != nil {
		name = document.Info.Name()
	} else {
		name = document.Path
	}

	return fmt.Sprintf("%s %s", name, document.Tags)
}
