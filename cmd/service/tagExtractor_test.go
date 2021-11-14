package service

import (
	"reflect"
	"testing"
)

func Test_extractTags(t *testing.T) {

	tests := []struct {
		name    string
		content string
		want    []string
	}{

		{
			name:    "no tags",
			content: "foo",
			want:    []string{},
		},

		{
			name:    "one tag",
			content: "#foo",
			want:    []string{"#foo"},
		},

		{
			name:    "two tags",
			content: "#tag0 #tag1",
			want:    []string{"#tag0", "#tag1"},
		},

		{
			name:    "no tag when markdown header",
			content: "# header content",
			want:    []string{},
		},

		{
			name: "multi line with markdown header and one tag",
			content: `
			# header
			content 
			#tag
			`,
			want: []string{"#tag"},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := ExtractTags(tt.content); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("ExtractTags() = %v, want %v", got, tt.want)
			}
		})
	}
}
