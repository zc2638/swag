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
