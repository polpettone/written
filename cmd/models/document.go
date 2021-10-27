package models

import (
	"fmt"
	"os"
)

type Document struct {
	Path string
	Info os.FileInfo
	Tags []string
}

func (document Document) String() string {
	return fmt.Sprintf("%s %s", document.Info.Name(), document.Tags)
}
