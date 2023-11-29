// Copyright © 2020 zc2638 <zc2638@qq.com>.
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
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNew(t *testing.T) {
	api := New()
	assert.Equal(t, "/", api.BasePath)
	assert.Equal(t, "2.0", api.Swagger)
	assert.Equal(t, []string{"http"}, api.Schemes)
	assert.Equal(t, []string{"http"}, api.Schemes)
	assert.Equal(t, "SNAPSHOT", api.Info.Version)
	assert.Equal(t, "https://swagger.io/terms/", api.Info.TermsOfService)
	assert.Equal(t, "Apache 2.0", api.Info.License.Name)
	assert.Equal(t, "https://www.apache.org/licenses/LICENSE-2.0.html", api.Info.License.URL)
}
