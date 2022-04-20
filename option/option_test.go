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

package option

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/zc2638/swag"
)

func TestDescription(t *testing.T) {
	api := swag.New(
		Description("zc"),
	)
	assert.Equal(t, "zc", api.Info.Description)
}

func TestVersion(t *testing.T) {
	api := swag.New(
		Version("zc"),
	)
	assert.Equal(t, "zc", api.Info.Version)
}

func TestTermsOfService(t *testing.T) {
	api := swag.New(
		TermsOfService("zc"),
	)
	assert.Equal(t, "zc", api.Info.TermsOfService)
}

func TestTitle(t *testing.T) {
	api := swag.New(
		Title("zc"),
	)
	assert.Equal(t, "zc", api.Info.Title)
}

func TestContactEmail(t *testing.T) {
	api := swag.New(
		ContactEmail("zc"),
	)
	assert.Equal(t, "zc", api.Info.Contact.Email)
}

func TestLicense(t *testing.T) {
	api := swag.New(
		License("name", "url"),
	)
	assert.Equal(t, "name", api.Info.License.Name)
	assert.Equal(t, "url", api.Info.License.URL)
}

func TestBasePath(t *testing.T) {
	api := swag.New(
		BasePath("/"),
	)
	assert.Equal(t, "/", api.BasePath)
}

func TestSchemes(t *testing.T) {
	api := swag.New(
		Schemes("zc"),
	)
	assert.Equal(t, []string{"zc"}, api.Schemes)
}

func TestHost(t *testing.T) {
	api := swag.New(
		Host("zc"),
	)
	assert.Equal(t, "zc", api.Host)
}

func TestSecurity(t *testing.T) {
	api := swag.New(
		Security("basic"),
	)
	assert.Len(t, api.Security.Requirements, 1)
	assert.Contains(t, api.Security.Requirements[0], "basic")
}
