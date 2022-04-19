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

import "fmt"

// New constructs a new api builder
func New(options ...Option) *API {
	api := &API{
		BasePath: "/",
		Swagger:  "2.0",
		Schemes:  []string{"http"},
		Info: Info{
			Description:    "Describe your API",
			Title:          "Your API Title",
			Version:        "SNAPSHOT",
			TermsOfService: "http://swagger.io/terms/",
			License: License{
				Name: "Apache 2.0",
				URL:  "http://www.apache.org/licenses/LICENSE-2.0.html",
			},
		},
	}

	for _, opt := range options {
		opt(api)
	}
	return api
}

// Option provides configuration options to the swagger api
type Option func(api *API)

// DescriptionOption sets info.description
func DescriptionOption(v string) Option {
	return func(api *API) {
		api.Info.Description = v
	}
}

// VersionOption sets info.version
func VersionOption(v string) Option {
	return func(api *API) {
		api.Info.Version = v
	}
}

// TermsOfServiceOption sets info.termsOfService
func TermsOfServiceOption(v string) Option {
	return func(api *API) {
		api.Info.TermsOfService = v
	}
}

// TitleOption sets info.title
func TitleOption(v string) Option {
	return func(api *API) {
		api.Info.Title = v
	}
}

// ContactEmailOption sets info.contact.email
func ContactEmailOption(v string) Option {
	return func(api *API) {
		if api.Info.Contact == nil {
			api.Info.Contact = &Contact{}
		}
		api.Info.Contact.Email = v
	}
}

// LicenseOption sets both info.license.name and info.license.url
func LicenseOption(name, url string) Option {
	return func(api *API) {
		api.Info.License.Name = name
		api.Info.License.URL = url
	}
}

// BasePathOption sets basePath
func BasePathOption(v string) Option {
	return func(api *API) {
		api.BasePath = v
	}
}

// SchemesOption sets the scheme
func SchemesOption(v ...string) Option {
	return func(api *API) {
		api.Schemes = v
	}
}

// TagOption adds a tag to the swagger api
func TagOption(tag *Tag) Option {
	return func(api *API) {
		api.Tags = append(api.Tags, *tag)
	}
}

// HostOption specifies the host field
func HostOption(v string) Option {
	return func(api *API) {
		api.Host = v
	}
}

// EndpointsOption allows the endpoints to be added dynamically to the Api
func EndpointsOption(endpoints ...*Endpoint) Option {
	return func(api *API) {
		api.AddEndpoint(endpoints...)
	}
}

// SecuritySchemeOption creates a new security definition for the API.
func SecuritySchemeOption(name string, scheme *SecurityScheme) Option {
	return func(api *API) {
		if api.SecurityDefinitions == nil {
			api.SecurityDefinitions = make(map[string]SecurityScheme)
		}
		api.SecurityDefinitions[name] = *scheme
	}
}

// SecurityOption sets a default security scheme for all endpoints in the API.
func SecurityOption(scheme string, scopes ...string) Option {
	return func(api *API) {
		if api.Security == nil {
			api.Security = &SecurityRequirement{}
		}

		if api.Security.Requirements == nil {
			api.Security.Requirements = []map[string][]string{}
		}
		if scopes == nil {
			scopes = make([]string, 0)
		}
		api.Security.Requirements = append(api.Security.Requirements, map[string][]string{scheme: scopes})
	}
}

// BasicSecurity defines a security scheme for HTTP Basic authentication.
func BasicSecurity() *SecurityScheme {
	return &SecurityScheme{Type: "basic"}
}

// APIKeySecurity defines a security scheme for API key authentication. "in" is
// the location of the API key (query or header). "name" is the name of the
// header or query parameter to be used.
func APIKeySecurity(name, in string) *SecurityScheme {
	if in != "header" && in != "query" {
		panic(fmt.Errorf(`APIKeySecurity "in" parameter must be one of: "header" or "query"`))
	}
	return &SecurityScheme{
		Type: "apiKey",
		Name: name,
		In:   in,
	}
}

// OAuth2Security defines a security scheme for OAuth2 authentication. Flow can
// be one of implicit, password, application, or accessCode.
func OAuth2Security(flow, authorizationURL, tokenURL string) *SecurityScheme {
	return &SecurityScheme{
		Type:             "oauth2",
		Flow:             flow,
		AuthorizationURL: authorizationURL,
		TokenURL:         tokenURL,
		Scopes:           make(map[string]string),
	}
}

// OAuth2ScopeSecurity adds a new scope to the security scheme.
//func OAuth2ScopeSecurity(scope, description string) *SecurityScheme {
//	return &SecurityScheme{
//		Scopes: map[string]string{
//			scope: description,
//		},
//	}
//}
