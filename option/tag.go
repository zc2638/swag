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

// Tag adds a tag to the swagger api
func Tag(name, description string, options ...TagOption) swag.Option {
	return func(api *swag.API) {
		t := swag.Tag{
			Name:        name,
			Description: description,
		}
		for _, opt := range options {
			opt(&t)
		}
		api.Tags = append(api.Tags, t)
	}
}

// TagOption provides additional customizations to the #Tag option
type TagOption func(tag *swag.Tag)

// TagDescription sets externalDocs.description on the tag field
func TagDescription(v string) TagOption {
	return func(t *swag.Tag) {
		if t.Docs == nil {
			t.Docs = &swag.TagDocs{}
		}
		t.Docs.Description = v
	}
}

// TagURL sets externalDocs.url on the tag field
func TagURL(v string) TagOption {
	return func(t *swag.Tag) {
		if t.Docs == nil {
			t.Docs = &swag.TagDocs{}
		}
		t.Docs.URL = v
	}
}
