package service

import (
	"github.com/polpettone/written/cmd/models"
	"io/ioutil"
	"os"
	"path/filepath"
)

func Read(path string) ([]*models.Document, error) {
	var documents []*models.Document

	if _, err := os.Stat(path); os.IsNotExist(err) {
		return nil, err
	}

	err := filepath.Walk(path,
		func(path string, info os.FileInfo, err error) error {

			if info.IsDir() &&  info.Name() == ".git" {
				return filepath.SkipDir
			}

			if info.IsDir() {
				return nil
			}

			content, err := ioutil.ReadFile(path)
			tags := ExtractTags(string(content))

			document := &models.Document{
				Path: path,
				Info: info,
				Tags: tags,
			}
			documents = append(documents, document)
			return nil
		})
	if err != nil {
		return nil, err
	}
	return documents, nil
}
