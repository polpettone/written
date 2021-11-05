package service

import (
	"fmt"
	"github.com/google/uuid"
	"github.com/polpettone/written/cmd/models"
	"io"
	"runtime"
	"strings"
	"testing"
)

func TestGetFilename(t *testing.T) {
	_, filename, _, _ := runtime.Caller(0)
	t.Logf("Current test filename: %s", filename)
}

func setupTestDocuments() []*models.Document {
	d0 := &models.Document{
		Path: "/tmp/d0",
		Info: nil,
		Tags: []string{"tag0", "tag1"},
	}
	d1 := &models.Document{
		Path: "/tmp/d1",
		Info: nil,
		Tags: []string{"tag0", "tag1"},
	}
	var documents []*models.Document
	documents = append(documents, d0)
	documents = append(documents, d1)
	return documents
}

func Test_Marshalling(t *testing.T) {
	documents := setupTestDocuments()
	r, err := Marshal(documents)

	if err != nil {
		t.Errorf("%v", err)
	}

	buf := new(strings.Builder)
	_, err = io.Copy(buf, r)

	jsonData := buf.String()

	fmt.Printf("%s", jsonData)

	var loaded []*models.Document
	err = Unmarshal(strings.NewReader(jsonData), &loaded)

	if err != nil {
		t.Errorf("%v", err)
	}

	if len(loaded) != 2 {
		t.Errorf("wanted %d got %d", 2, len(loaded))
	}
}

func Test_Read_Documents(t *testing.T) {
	documents, err := Read("../../test-data/documents")
	if err != nil {
		t.Errorf("%v", err)
	}
	if len(documents) != 2 {
		t.Errorf("wanted %d got %d", 1, len(documents))
	}
}

func Test_UpdateDocuments(t *testing.T) {
	id, _ := uuid.NewUUID()
	path := "/tmp/written-test-documents-" + id.String()

	documents,_ := Read("../../test-data/documents")
	documents[0].Tags = append(documents[0].Tags, "neu")

	_ = Save(path, documents)
	loaded, _ := Load(path, documents)

	if loaded[0].Info == nil {
		t.Errorf("File Info should not be nil")
	}
}

func Test_SaveAndLoad(t *testing.T) {
	id, _ := uuid.NewUUID()
	path := "/tmp/written-test-documents-" + id.String()
	documents, _ := Read("../../test-data/documents")

	err := Save(path, documents)
	if err != nil {
		t.Errorf("%v", err)
	}

	loaded, err := Load(path, documents)

	if err != nil {
		t.Errorf("%v", err)
	}

	if len(loaded) != 2 {
		t.Errorf("wanted %d got %d", 2, len(loaded))
	}
}
