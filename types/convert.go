// Copyright © 2024 zc2638 <zc2638@qq.com>.
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

package types

import (
	"reflect"
)

var ptMapping = map[string]ParameterType{
	"time.Time":     "string",
	"time.Duration": "integer",
	"json.Number":   "number",
}

func Register(goType string, pt ParameterType) {
	ptMapping[goType] = pt
}

func Get(goType string) ParameterType {
	return ptMapping[goType]
}

func Parse(t reflect.Type) ParameterType {
	kind := t.Kind()
	if kind == reflect.Ptr {
		return Parse(t.Elem())
	}

	if pt := Get(t.String()); pt != Unknown {
		return pt
	}

	switch kind {
	case reflect.Bool:
		return Boolean
	case reflect.Float32, reflect.Float64:
		return Number
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64,
		reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		return Integer
	case reflect.String:
		return String
	case reflect.Array, reflect.Slice:
		return Array
	case reflect.Struct, reflect.Map:
		return Object
	default:
		return Unknown
	}
}
