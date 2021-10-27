package service

import (
	"github.com/polpettone/written/cmd/config"
	"github.com/polpettone/written/cmd/models"
	"os"
	"path/filepath"
)

func Read(path string) ([]*models.Document, error) {
	var documents []*models.Document
	err := filepath.Walk(path,
		func(path string, info os.FileInfo, err error) error {

			if info.IsDir() &&  info.Name() == ".git" {
				return filepath.SkipDir
			}

			if info.IsDir() {
				return nil
			}

			document := &models.Document{
				Path: path,
				Info: info,
				Tags: []string{"dummy0", "dummy1"},
			}
			documents = append(documents, document)
			return nil
		})
	if err != nil {
		return nil, err
	}
	return documents, nil
}

func Save() {
	config.Log.InfoLog.Printf("Save")
}
