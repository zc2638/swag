// Created by zc on 2022/10/31.

package swag

import (
	"fmt"
	"io/fs"
	"net/http"
	"path"
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
				dir:  asserts.DistDir,
				fsys: asserts.Dist,
			},
			want: dirFS{
				dir: asserts.DistDir,
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
	file, err := http.FS(asserts.Dist).Open(path.Join(asserts.DistDir, "index.html"))
	if err != nil {
		t.Errorf("open asserts failed: %v", err)
	}

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
		{
			name: "",
			fields: fields{
				dir: asserts.DistDir,
				fs:  http.FS(asserts.Dist),
			},
			args: args{
				name: "index.html",
			},
			want:    file,
			wantErr: func(t assert.TestingT, err error, i ...interface{}) bool { return err == nil },
		},
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
