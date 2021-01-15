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
//
package swagger

import (
	"fmt"
	"github.com/modern-go/reflect2"
	"path/filepath"
	"reflect"
	"strconv"
	"strings"
)

func makeRef(name string) string {
	return fmt.Sprintf("#/definitions/%v", name)
}

func makeName(t reflect.Type) string {
	name := t.Name()
	if name == "" {
		ptr := reflect2.PtrOf(t)
		name = "ptr" + strconv.FormatUint(uint64(uintptr(ptr)), 10)
	}
	full := filepath.Base(t.PkgPath()) + name
	return strings.Replace(full, "-", "_", -1)
}
