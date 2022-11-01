// Created by zc on 2022/11/1.

package swag

import (
	"fmt"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestEndpoint_BuildOperationID(t *testing.T) {
	type fields struct {
		Path   string
		Method string
	}
	tests := []struct {
		name   string
		fields fields
		want   string
	}{
		{
			name: "case 1: none",
			want: "",
		},
		{
			name: "case 2: get",
			fields: fields{
				Path:   "/",
				Method: http.MethodGet,
			},
			want: "get",
		},
		{
			name: "case 3: post",
			fields: fields{
				Path:   "/test",
				Method: http.MethodPost,
			},
			want: "postTest",
		},
		{
			name: "case 4: put",
			fields: fields{
				Path:   "/test/{id}",
				Method: http.MethodPut,
			},
			want: "putTestId",
		},
		{
			name: "case 5: patch",
			fields: fields{
				Path:   "/test/{id}/sub",
				Method: http.MethodPatch,
			},
			want: "patchTestIdSub",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			e := &Endpoint{
				Path:   tt.fields.Path,
				Method: tt.fields.Method,
			}
			e.BuildOperationID()
			assert.Equal(t, tt.want, e.OperationID)
		})
	}
}

func TestSecurityRequirement_MarshalJSON(t *testing.T) {
	nilErrorFunc := func(t assert.TestingT, err error, i ...interface{}) bool {
		return err == nil
	}

	type fields struct {
		Requirements    []map[string][]string
		DisableSecurity bool
	}
	tests := []struct {
		name    string
		fields  fields
		want    []byte
		wantErr assert.ErrorAssertionFunc
	}{
		{
			name:    "case 1: null",
			fields:  fields{},
			want:    []byte("null"),
			wantErr: nilErrorFunc,
		},
		{
			name: "case 2: disable security",
			fields: fields{
				Requirements: []map[string][]string{
					{"oauth": []string{"scope1", "scope2"}},
				},
				DisableSecurity: true,
			},
			want:    []byte("[]"),
			wantErr: nilErrorFunc,
		},
		{
			name: "case 3: disable security without requirements",
			fields: fields{
				Requirements:    nil,
				DisableSecurity: true,
			},
			want:    []byte("[]"),
			wantErr: nilErrorFunc,
		},
		{
			name: "case 4: enable security",
			fields: fields{
				Requirements: []map[string][]string{
					{"oauth": []string{"scope1", "scope2"}},
				},
				DisableSecurity: false,
			},
			want:    []byte(`[{"oauth":["scope1","scope2"]}]`),
			wantErr: nilErrorFunc,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &SecurityRequirement{
				Requirements:    tt.fields.Requirements,
				DisableSecurity: tt.fields.DisableSecurity,
			}
			got, err := s.MarshalJSON()
			if !tt.wantErr(t, err, fmt.Sprintf("MarshalJSON()")) {
				return
			}
			assert.Equalf(t, string(tt.want), string(got), "MarshalJSON()")
		})
	}
}
