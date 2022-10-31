// Copyright © 2022 zc2638 <zc2638@qq.com>.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package swag

import (
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestEndpoints_ServeHTTPNotFound(t *testing.T) {
	req, _ := http.NewRequest(http.MethodGet, "http://localhost", nil)
	w := httptest.NewRecorder()

	var es Endpoints
	es.ServeHTTP(w, req)
	assert.Equal(t, http.StatusNotFound, w.Code)
}

func TestEndpoints_ServeHTTP(t *testing.T) {
	fn := func(v string) *Endpoint {
		return &Endpoint{
			Handler: func(w http.ResponseWriter, req *http.Request) {
				w.WriteHeader(http.StatusOK)
				_, _ = io.WriteString(w, v)
			},
		}
	}

	es := Endpoints{
		Delete:  fn("Delete"),
		Head:    fn("Head"),
		Get:     fn("Get"),
		Options: fn("Options"),
		Post:    fn("Post"),
		Put:     fn("Put"),
		Patch:   fn("Patch"),
		Trace:   fn("Trace"),
		Connect: fn("Connect"),
	}

	methods := []string{
		http.MethodDelete,
		http.MethodHead,
		http.MethodGet,
		http.MethodOptions,
		http.MethodPost,
		http.MethodPut,
		http.MethodPatch,
		http.MethodTrace,
		http.MethodConnect,
	}
	for _, method := range methods {
		req, err := http.NewRequest(strings.ToUpper(method), "http://localhost", nil)
		assert.Nil(t, err)

		w := httptest.NewRecorder()
		es.ServeHTTP(w, req)
		assert.Equal(t, strings.ToUpper(w.Body.String()), strings.ToUpper(method))
	}
}

func TestAPI_AddOptions(t *testing.T) {
	type args struct {
		options []Option
	}
	tests := []struct {
		name string
		args args
		want *API
	}{
		{
			name: "base path",
			args: args{
				options: []Option{
					func(api *API) {
						api.BasePath = "/test"
					},
				},
			},
			want: &API{
				BasePath: "/test",
				Swagger:  "2.0",
				Schemes:  []string{"http"},
				Info: Info{
					Description:    "Describe your API",
					Title:          "Your API Title",
					Version:        "SNAPSHOT",
					TermsOfService: "http://swagger.io/terms/",
					License: License{
						Name: "Apache 2.0",
						URL:  "http://www.apache.org/licenses/LICENSE-2.0.html",
					},
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			actual := New()
			actual.AddOptions(tt.args.options...)
			assert.Equal(t, tt.want, actual)
		})
	}
}

func TestAPI_AddTag(t *testing.T) {
	type args struct {
		name        string
		description string
	}
	tests := []struct {
		name string
		args args
		want *API
	}{
		{
			name: "add test tag",
			args: args{
				name:        "test",
				description: "test",
			},
			want: &API{
				BasePath: "/",
				Swagger:  "2.0",
				Schemes:  []string{"http"},
				Info: Info{
					Description:    "Describe your API",
					Title:          "Your API Title",
					Version:        "SNAPSHOT",
					TermsOfService: "http://swagger.io/terms/",
					License: License{
						Name: "Apache 2.0",
						URL:  "http://www.apache.org/licenses/LICENSE-2.0.html",
					},
				},
				Tags: []Tag{
					{
						Name:        "test",
						Description: "test",
					},
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			actual := New()
			actual.AddTag(tt.args.name, tt.args.description)
			assert.Equal(t, tt.want, actual)
		})
	}
}

func TestAPI_Clone(t *testing.T) {
	type fields struct {
		Swagger             string
		Info                Info
		BasePath            string
		Schemes             []string
		Paths               map[string]*Endpoints
		Definitions         map[string]Object
		Tags                []Tag
		Host                string
		SecurityDefinitions map[string]SecurityScheme
		Security            *SecurityRequirement
		tags                []Tag
		prefixPath          string
	}
	tests := []struct {
		name   string
		fields fields
		want   *API
	}{
		{
			name:   "none",
			fields: fields{},
			want:   &API{},
		},
		{
			name: "base field",
			fields: fields{
				BasePath: "/",
				Swagger:  "2.0",
				Schemes:  []string{"http"},
			},
			want: &API{
				BasePath: "/",
				Swagger:  "2.0",
				Schemes:  []string{"http"},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			a := &API{
				Swagger:             tt.fields.Swagger,
				Info:                tt.fields.Info,
				BasePath:            tt.fields.BasePath,
				Schemes:             tt.fields.Schemes,
				Paths:               tt.fields.Paths,
				Definitions:         tt.fields.Definitions,
				Tags:                tt.fields.Tags,
				Host:                tt.fields.Host,
				SecurityDefinitions: tt.fields.SecurityDefinitions,
				Security:            tt.fields.Security,
				tags:                tt.fields.tags,
				prefixPath:          tt.fields.prefixPath,
			}
			assert.Equalf(t, tt.want, a.Clone(), "Clone()")
		})
	}
}

func TestUIHandler(t *testing.T) {
	type args struct {
		prefix string
		uri    string
		req    *http.Request
	}
	tests := []struct {
		name     string
		args     args
		wantCode int
	}{
		{
			name: "200",
			args: args{
				prefix: "/swagger/ui",
				uri:    "/swagger/ui",
				req:    httptest.NewRequest(http.MethodGet, "/swagger/ui/", nil),
			},
			wantCode: http.StatusOK,
		},
		{
			name: "302",
			args: args{
				prefix: "/swagger/ui",
				uri:    "/swagger/ui",
				req:    httptest.NewRequest(http.MethodGet, "/swagger/ui", nil),
			},
			wantCode: http.StatusFound,
		},
		{
			name: "404",
			args: args{
				prefix: "/swagger/ui",
				uri:    "/swagger/ui",
				req:    httptest.NewRequest(http.MethodGet, "/", nil),
			},
			wantCode: http.StatusNotFound,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			w := httptest.NewRecorder()
			handler := UIHandler(tt.args.prefix, tt.args.uri, false)
			handler.ServeHTTP(w, tt.args.req)
			assert.Equalf(t, tt.wantCode, w.Code, "UIPatterns(%v)", tt.args.prefix)
		})
	}
}

func TestUIPatterns(t *testing.T) {
	all := []string{
		"/",
		"favicon-16x16.png",
		"favicon-32x32.png",
		"index.html",
		"oauth2-redirect.html",
		"swagger-ui-bundle.js",
		"swagger-ui-bundle.js.map",
		"swagger-ui-es-bundle-core.js",
		"swagger-ui-es-bundle-core.js.map",
		"swagger-ui-es-bundle.js",
		"swagger-ui-es-bundle.js.map",
		"swagger-ui-standalone-preset.js",
		"swagger-ui-standalone-preset.js.map",
		"swagger-ui.css",
		"swagger-ui.css.map",
		"swagger-ui.js",
		"swagger-ui.js.map",
	}
	uiAll := []string{
		"/swagger/ui/",
		"/swagger/ui/favicon-16x16.png",
		"/swagger/ui/favicon-32x32.png",
		"/swagger/ui/index.html",
		"/swagger/ui/oauth2-redirect.html",
		"/swagger/ui/swagger-ui-bundle.js",
		"/swagger/ui/swagger-ui-bundle.js.map",
		"/swagger/ui/swagger-ui-es-bundle-core.js",
		"/swagger/ui/swagger-ui-es-bundle-core.js.map",
		"/swagger/ui/swagger-ui-es-bundle.js",
		"/swagger/ui/swagger-ui-es-bundle.js.map",
		"/swagger/ui/swagger-ui-standalone-preset.js",
		"/swagger/ui/swagger-ui-standalone-preset.js.map",
		"/swagger/ui/swagger-ui.css",
		"/swagger/ui/swagger-ui.css.map",
		"/swagger/ui/swagger-ui.js",
		"/swagger/ui/swagger-ui.js.map",
	}

	type args struct {
		prefix string
	}
	tests := []struct {
		name string
		args args
		want []string
	}{
		{
			name: "no prefix",
			args: args{prefix: ""},
			want: all,
		},
		{
			name: "swagger prefix",
			args: args{prefix: "/swagger/ui"},
			want: uiAll,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equalf(t, tt.want, UIPatterns(tt.args.prefix), "UIPatterns(%v)", tt.args.prefix)
		})
	}
}
