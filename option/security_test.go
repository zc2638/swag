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

func TestSecurityScheme(t *testing.T) {
	api := swag.New(
		SecurityScheme("basic", BasicSecurity()),
		SecurityScheme("apikey", APIKeySecurity("Authorization", "header")),
	)
	assert.Len(t, api.SecurityDefinitions, 2)
	assert.Contains(t, api.SecurityDefinitions, "basic")
	assert.Contains(t, api.SecurityDefinitions, "apikey")
	assert.Equal(t, "header", api.SecurityDefinitions["apikey"].In)
}

func TestSecuritySchemeDescription(t *testing.T) {
	scheme := &swag.SecurityScheme{}
	description := "a security scheme"
	SecuritySchemeDescription(description)(scheme)
	assert.Equal(t, description, scheme.Description)
}

func TestBasicSecurity(t *testing.T) {
	scheme := &swag.SecurityScheme{}
	BasicSecurity()(scheme)
	assert.Equal(t, scheme.Type, "basic")
}

func TestAPIKeySecurity(t *testing.T) {
	scheme := &swag.SecurityScheme{}
	name := "Authorization"
	in := "header"
	APIKeySecurity(name, in)(scheme)
	assert.Equal(t, scheme.Type, "apiKey")
	assert.Equal(t, scheme.Name, name)
	assert.Equal(t, scheme.In, in)

	assert.Panics(t,
		func() { APIKeySecurity(name, "invalid") },
		"expected APIKeySecurity to panic with invalid \"in\" parameter",
	)
}

func TestOAuth2Security(t *testing.T) {
	scheme := &swag.SecurityScheme{}
	flow := "accessCode"
	authURL := "https://example.com/oauth/authorize"
	tokenURL := "https://example.com/oauth/token"
	OAuth2Security(flow, authURL, tokenURL)(scheme)

	assert.Equal(t, scheme.Type, "oauth2")
	assert.Equal(t, scheme.Flow, "accessCode")
	assert.Equal(t, scheme.AuthorizationURL, authURL)
	assert.Equal(t, scheme.TokenURL, tokenURL)
}

func TestOAuth2Scope(t *testing.T) {
	scheme := &swag.SecurityScheme{}

	OAuth2Scope("read", "read data")(scheme)
	OAuth2Scope("write", "write data")(scheme)

	assert.Len(t, scheme.Scopes, 2)
	assert.Contains(t, scheme.Scopes, "read")
	assert.Contains(t, scheme.Scopes, "write")

	assert.Equal(t, "read data", scheme.Scopes["read"])
	assert.Equal(t, "write data", scheme.Scopes["write"])
}
