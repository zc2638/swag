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

import "reflect"

type Format string

func (f Format) String() string {
	return string(f)
}

const (
	FormatNone   Format = ""
	FormatInt32  Format = "int32"
	FormatInt64  Format = "int64"
	FormatDouble Format = "double"
	FormatFloat  Format = "float"
)

func ParseFormat(t reflect.Type) Format {
	kind := t.Kind()
	if kind == reflect.Ptr {
		return ParseFormat(t.Elem())
	}

	switch kind {
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Uint8, reflect.Uint16, reflect.Uint32:
		return FormatInt32
	case reflect.Int64, reflect.Uint64:
		return FormatInt64
	case reflect.Float64:
		return FormatDouble
	case reflect.Float32:
		return FormatFloat
	default:
		return FormatNone
	}
}
