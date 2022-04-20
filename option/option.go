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
	"github.com/zc2638/swag"
)

// Description sets info.description
func Description(v string) swag.Option {
	return func(api *swag.API) {
		api.Info.Description = v
	}
}

// Version sets info.version
func Version(v string) swag.Option {
	return func(api *swag.API) {
		api.Info.Version = v
	}
}

// TermsOfService sets info.termsOfService
func TermsOfService(v string) swag.Option {
	return func(api *swag.API) {
		api.Info.TermsOfService = v
	}
}

// Title sets info.title
func Title(v string) swag.Option {
	return func(api *swag.API) {
		api.Info.Title = v
	}
}

// ContactEmail sets info.contact.email
func ContactEmail(v string) swag.Option {
	return func(api *swag.API) {
		if api.Info.Contact == nil {
			api.Info.Contact = &swag.Contact{}
		}
		api.Info.Contact.Email = v
	}
}

// License sets both info.license.name and info.license.url
func License(name, url string) swag.Option {
	return func(api *swag.API) {
		api.Info.License.Name = name
		api.Info.License.URL = url
	}
}

// BasePath sets basePath
func BasePath(v string) swag.Option {
	return func(api *swag.API) {
		api.BasePath = v
	}
}

// Schemes sets the scheme
func Schemes(v ...string) swag.Option {
	return func(api *swag.API) {
		api.Schemes = v
	}
}

// Host specifies the host field
func Host(v string) swag.Option {
	return func(api *swag.API) {
		api.Host = v
	}
}

// Endpoints allows the endpoints to be added dynamically to the Api
func Endpoints(endpoints ...*swag.Endpoint) swag.Option {
	return func(api *swag.API) {
		api.AddEndpoint(endpoints...)
	}
}

// Security sets a default security scheme for all endpoints in the API.
func Security(scheme string, scopes ...string) swag.Option {
	return func(api *swag.API) {
		if api.Security == nil {
			api.Security = &swag.SecurityRequirement{}
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
