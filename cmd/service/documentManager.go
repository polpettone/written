package service

import (
	"bytes"
	"encoding/json"
	"github.com/polpettone/written/cmd/config"
	"github.com/polpettone/written/cmd/models"
	"io"
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

			document := &models.Document{
				Path: path,
				Info: info,
				Tags: []string{},
			}
			documents = append(documents, document)
			return nil
		})
	if err != nil {
		return nil, err
	}
	return documents, nil
}

var Unmarshal = func(r io.Reader, v interface{}) error {
	return json.NewDecoder(r).Decode(v)
}

var Marshal = func(v interface{}) (io.Reader, error) {
	b, err := json.MarshalIndent(v, "", "\t")
	if err != nil {
		return nil, err
	}
	return bytes.NewReader(b), nil
}

func Save(path string, documents []*models.Document) error {
	jsonData, err := json.Marshal(documents)
	if err != nil {
		return err
	}
	err = ioutil.WriteFile(path, jsonData, 0644)
	if err != nil {
		return err
	}
	config.Log.InfoLog.Printf("Saved documents to %s", path)
	return nil
}

func Load(path string, currentDocuments []*models.Document) ([]*models.Document, error) {
	f, err := os.Open(path)
	if err != nil {
		return currentDocuments, err
	}
	var loaded []*models.Document
	err = Unmarshal(f, &loaded)
	if err != nil {
		return nil, err
	}

	for _, currentDocument := range currentDocuments {
		for _, document := range loaded {
			if currentDocument.Path == document.Path {
				currentDocument.Tags = append(currentDocument.Tags, document.Tags...)
			}
		}
	}

	return currentDocuments, err
}

