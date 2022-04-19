// Copyright 2020 zc2638
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
// http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package swag

import (
	"reflect"
	"strconv"
	"testing"

	"github.com/modern-go/reflect2"

	"github.com/stretchr/testify/assert"
)

func TestColonPath(t *testing.T) {
	assert.Equal(t, "/api/:id", ColonPath("/api/{id}"))
	assert.Equal(t, "/api/:a/:b/:c", ColonPath("/api/{a}/{b}/{c}"))
}

func getPtrString(t reflect.Type) string {
	ptr := reflect2.PtrOf(t)
	return ".ptr" + strconv.FormatUint(uint64(uintptr(ptr)), 10)
}

func Test_makeName(t *testing.T) {
	test := struct {
		ID   string `json:"id"`
		Name string `json:"name"`
	}{}
	test2 := struct {
		ID   string `json:"id"`
		Data string `json:"data"`
	}{}

	type args struct {
		t reflect.Type
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			name: "test",
			args: args{t: reflect.TypeOf(test)},
			want: getPtrString(reflect.TypeOf(test)),
		},
		{
			name: "test2",
			args: args{t: reflect.TypeOf(test2)},
			want: getPtrString(reflect.TypeOf(test2)),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := makeName(tt.args.t); got != tt.want {
				t.Errorf("makeName() = %v, want %v", got, tt.want)
			}
		})
	}
}
