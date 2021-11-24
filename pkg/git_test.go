package pkg

import "testing"

func Test_getRootPathOfFile(t *testing.T) {
	type args struct {
		rootPath string
		filePath string
	}
	tests := []struct {
		name string
		args args
		want string
	}{

		{
			name: "one",
			args: args{
				rootPath: "root",
				filePath: "root/file",
			},
			want: "file",
		},
		{
			name: "two",
			args: args{
				rootPath: "/home/akim/storage/written-documents",
				filePath: "/home/akim/storage/written-documents/nlc/Secret-nlc",
			},
			want: "nlc/Secret-nlc",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := getRelPathToRootPath(tt.args.rootPath, tt.args.filePath); got != tt.want {
				t.Errorf("getRootPathOfFile() = %v, want %v", got, tt.want)
			}
		})
	}
}
