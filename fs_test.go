// Created by zc on 2022/10/31.

package swag

import (
	"fmt"
	"io/fs"
	"net/http"
	"testing"

	"github.com/zc2638/swag/asserts"

	"github.com/stretchr/testify/assert"
)

func TestDirFS(t *testing.T) {
	type args struct {
		dir  string
		fsys fs.FS
	}
	tests := []struct {
		name string
		args args
		want http.FileSystem
	}{
		{
			name: "embed",
			args: args{
				dir:  "/swagger",
				fsys: asserts.Dist,
			},
			want: dirFS{
				dir: "/swagger",
				fs:  http.FS(asserts.Dist),
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equalf(t, tt.want, DirFS(tt.args.dir, tt.args.fsys), "DirFS(%v, %v)", tt.args.dir, tt.args.fsys)
		})
	}
}

func Test_dirFS_Open(t *testing.T) {
	type fields struct {
		dir string
		fs  http.FileSystem
	}
	type args struct {
		name string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    http.File
		wantErr assert.ErrorAssertionFunc
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			f := dirFS{
				dir: tt.fields.dir,
				fs:  tt.fields.fs,
			}
			got, err := f.Open(tt.args.name)
			if !tt.wantErr(t, err, fmt.Sprintf("Open(%v)", tt.args.name)) {
				return
			}
			assert.Equalf(t, tt.want, got, "Open(%v)", tt.args.name)
		})
	}
}
