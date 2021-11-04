package models

import (
	"fmt"
	"os"
)

type Document struct {
	Path string `json:"Path"`
	Info os.FileInfo `json:"-"`
	Tags []string `json:"Tags"`
}

func (document Document) String() string {
	return fmt.Sprintf("%s %s", document.Info.Name(), document.Tags)
}
