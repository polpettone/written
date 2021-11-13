package service

import (
	"runtime"
	"testing"
)

func TestGetFilename(t *testing.T) {
	_, filename, _, _ := runtime.Caller(0)
	t.Logf("Current test filename: %s", filename)
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
