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
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestDescription(t *testing.T) {
	api := New(
		DescriptionOption("zc"),
	)
	assert.Equal(t, "zc", api.Info.Description)
}

func TestVersion(t *testing.T) {
	api := New(
		VersionOption("zc"),
	)
	assert.Equal(t, "zc", api.Info.Version)
}

func TestTermsOfService(t *testing.T) {
	api := New(
		TermsOfServiceOption("zc"),
	)
	assert.Equal(t, "zc", api.Info.TermsOfService)
}

func TestTitle(t *testing.T) {
	api := New(
		TitleOption("zc"),
	)
	assert.Equal(t, "zc", api.Info.Title)
}

func TestContactEmail(t *testing.T) {
	api := New(
		ContactEmailOption("zc"),
	)
	assert.Equal(t, "zc", api.Info.Contact.Email)
}

func TestLicense(t *testing.T) {
	api := New(
		LicenseOption("name", "url"),
	)
	assert.Equal(t, "name", api.Info.License.Name)
	assert.Equal(t, "url", api.Info.License.URL)
}

func TestBasePath(t *testing.T) {
	api := New(
		BasePathOption("/"),
	)
	assert.Equal(t, "/", api.BasePath)
}

func TestSchemes(t *testing.T) {
	api := New(
		SchemesOption("zc"),
	)
	assert.Equal(t, []string{"zc"}, api.Schemes)
}

func TestTag(t *testing.T) {
	api := New(
		TagOption(&Tag{
			Name:        "name",
			Description: "desc",
			Docs: &TagDocs{
				Description: "ext-desc",
				URL:         "ext-url",
			},
		}),
	)

	expected := Tag{
		Name:        "name",
		Description: "desc",
		Docs: &TagDocs{
			Description: "ext-desc",
			URL:         "ext-url",
		},
	}
	assert.Equal(t, expected, api.Tags[0])
}

func TestHost(t *testing.T) {
	api := New(
		HostOption("zc"),
	)
	assert.Equal(t, "zc", api.Host)
}

func TestSecurityScheme(t *testing.T) {
	api := New(
		SecuritySchemeOption("basic", BasicSecurity()),
		SecuritySchemeOption("apikey", APIKeySecurity("Authorization", "header")),
	)
	assert.Len(t, api.SecurityDefinitions, 2)
	assert.Contains(t, api.SecurityDefinitions, "basic")
	assert.Contains(t, api.SecurityDefinitions, "apikey")
	assert.Equal(t, "header", api.SecurityDefinitions["apikey"].In)
}

func TestSecurity(t *testing.T) {
	api := New(
		SecurityOption("basic"),
	)
	assert.Len(t, api.Security.Requirements, 1)
	assert.Contains(t, api.Security.Requirements[0], "basic")
}

func TestBasicSecurity(t *testing.T) {
	scheme := BasicSecurity()
	assert.Equal(t, scheme.Type, "basic")
}

func TestAPIKeySecurity(t *testing.T) {
	name := "Authorization"
	in := "header"
	scheme := APIKeySecurity(name, in)
	assert.Equal(t, scheme.Type, "apiKey")
	assert.Equal(t, scheme.Name, name)
	assert.Equal(t, scheme.In, in)

	assert.Panics(t,
		func() { APIKeySecurity(name, "invalid") },
		"expected APIKeySecurity to panic with invalid \"in\" parameter",
	)
}

func TestOAuth2Security(t *testing.T) {
	flow := "accessCode"
	authURL := "https://example.com/oauth/authorize"
	tokenURL := "https://example.com/oauth/token"
	scheme := OAuth2Security(flow, authURL, tokenURL)

	assert.Equal(t, scheme.Type, "oauth2")
	assert.Equal(t, scheme.Flow, "accessCode")
	assert.Equal(t, scheme.AuthorizationURL, authURL)
	assert.Equal(t, scheme.TokenURL, tokenURL)
}
