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
