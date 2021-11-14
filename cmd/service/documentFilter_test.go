package service

import (
	"github.com/polpettone/written/cmd/models"
	"reflect"
	"testing"
)

func TestFilterDocuments(t *testing.T) {
	type args struct {
		documents []*models.Document
		query     string
	}
	tests := []struct {
		name string
		args args
		want []*models.Document
	}{

		{
			name: "TestCase 0",
			args: args{
				documents: []*models.Document{
					createDocument("0", "foobar"),
					createDocument("1", "content")},
				query: "foobar",
			},
			want: []*models.Document{
				createDocument("0", "foobar")},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := FilterDocuments(tt.args.documents, tt.args.query); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("FilterDocuments() = %v, want %v", got, tt.want)
			}
		})
	}
}

func createDocument(path string, content string) *models.Document {
	return &models.Document{
		Content: content,
		Path:    path,
	}
}
