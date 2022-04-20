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
	"fmt"

	"github.com/zc2638/swag"
)

// SecurityScheme creates a new security definition for the API.
func SecurityScheme(name string, options ...SecuritySchemeOption) swag.Option {
	return func(api *swag.API) {
		if api.SecurityDefinitions == nil {
			api.SecurityDefinitions = make(map[string]swag.SecurityScheme)
		}

		var scheme swag.SecurityScheme
		for _, opt := range options {
			opt(&scheme)
		}
		api.SecurityDefinitions[name] = scheme
	}
}

// SecuritySchemeOption provides additional customizations to the SecurityScheme.
type SecuritySchemeOption func(securityScheme *swag.SecurityScheme)

// SecuritySchemeDescription sets the security scheme's description.
func SecuritySchemeDescription(description string) SecuritySchemeOption {
	return func(securityScheme *swag.SecurityScheme) {
		securityScheme.Description = description
	}
}

// BasicSecurity defines a security scheme for HTTP Basic authentication.
func BasicSecurity() SecuritySchemeOption {
	return func(securityScheme *swag.SecurityScheme) {
		securityScheme.Type = "basic"
	}
}

// APIKeySecurity defines a security scheme for API key authentication. "in" is
// the location of the API key (query or header). "name" is the name of the
// header or query parameter to be used.
func APIKeySecurity(name, in string) SecuritySchemeOption {
	if in != "header" && in != "query" {
		panic(fmt.Errorf(`APIKeySecurity "in" parameter must be one of: "header" or "query"`))
	}
	return func(scheme *swag.SecurityScheme) {
		scheme.Type = "apiKey"
		scheme.Name = name
		scheme.In = in
	}
}

// OAuth2Scope adds a new scope to the security scheme.
func OAuth2Scope(scope, description string) SecuritySchemeOption {
	return func(scheme *swag.SecurityScheme) {
		if scheme.Scopes == nil {
			scheme.Scopes = make(map[string]string)
		}
		scheme.Scopes[scope] = description
	}
}

// OAuth2Security defines a security scheme for OAuth2 authentication. Flow can
// be one of implicit, password, application, or accessCode.
func OAuth2Security(flow, authorizationURL, tokenURL string) SecuritySchemeOption {
	return func(scheme *swag.SecurityScheme) {
		scheme.Type = "oauth2"
		scheme.Flow = flow
		scheme.AuthorizationURL = authorizationURL
		scheme.TokenURL = tokenURL
		if scheme.Scopes == nil {
			scheme.Scopes = make(map[string]string)
		}
	}
}
