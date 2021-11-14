package service

import (
	"github.com/polpettone/written/cmd/models"
	"strings"
)

func FilterDocuments(documents []*models.Document, query string) []*models.Document {
	var filteredDocuments []*models.Document
	if query != "" {
		for _, d := range documents {
			if strings.Contains(d.Content, query) {
				filteredDocuments = append(filteredDocuments, d)
			}
		}
	} else {
		filteredDocuments = documents
	}
	return filteredDocuments
}