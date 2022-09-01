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

package endpoint

import (
	"io"
	"net/http"
	"reflect"
	"testing"

	"github.com/zc2638/swag/types"

	"github.com/zc2638/swag/option"

	"github.com/stretchr/testify/assert"

	"github.com/zc2638/swag"
)

func Echo(w http.ResponseWriter, _ *http.Request) {
	_, _ = io.WriteString(w, "Hello World!")
}

func TestNew(t *testing.T) {
	summary := "here's the summary"
	e := New(
		"get", "/",
		Summary(summary),
		Handler(Echo),
	)

	assert.Equal(t, "GET", e.Method)
	assert.Equal(t, "/", e.Path)
	assert.NotNil(t, e.Handler)
	assert.Equal(t, []string{"application/json"}, e.Consumes)
	assert.Equal(t, []string{"application/json"}, e.Produces)
	assert.Equal(t, summary, e.Summary)
	assert.Equal(t, []string(nil), e.Tags)
}

func TestTags(t *testing.T) {
	e := New(
		"get", "/",
		Summary("get thing"),
		Tags("zc"),
	)
	assert.Equal(t, []string{"zc"}, e.Tags)
}

func TestDescription(t *testing.T) {
	e := New(
		"get", "/",
		Summary("get thing"),
		Description("zc"),
	)

	assert.Equal(t, "zc", e.Description)
}

func TestOperationId(t *testing.T) {
	e := New(
		"get", "/",
		Summary("get thing"),
		OperationID("zc"),
	)

	assert.Equal(t, "zc", e.OperationID)
}

func TestProduces(t *testing.T) {
	expected := []string{"a", "b"}
	e := New(
		"get", "/",
		Summary("get thing"),
		Produces(expected...),
	)

	assert.Equal(t, expected, e.Produces)
}

func TestConsumes(t *testing.T) {
	expected := []string{"a", "b"}
	e := New(
		"get", "/",
		Summary("get thing"),
		Consumes(expected...),
	)

	assert.Equal(t, expected, e.Consumes)
}

func TestPath(t *testing.T) {
	expected := swag.Parameter{
		In:          "path",
		Name:        "id",
		Description: "the description",
		Required:    true,
		Type:        "string",
	}

	e := New(
		"get", "/",
		Summary("get thing"),
		Path(expected.Name, expected.Type, expected.Description, expected.Required),
	)

	assert.Equal(t, 1, len(e.Parameters))
	assert.Equal(t, expected, e.Parameters[0])
}

func TestPathString(t *testing.T) {
	expected := swag.Parameter{
		In:          "path",
		Name:        "id",
		Description: "the description",
		Required:    true,
		Type:        "string",
	}

	e := New(
		"get", "/",
		Summary("get thing"),
		PathS(expected.Name, expected.Description),
	)

	assert.Equal(t, 1, len(e.Parameters))
	assert.Equal(t, expected, e.Parameters[0])
}

func TestQuery(t *testing.T) {
	expected := swag.Parameter{
		In:          "query",
		Name:        "id",
		Description: "the description",
		Required:    true,
		Type:        "string",
	}

	e := New("get", "/",
		Summary("get thing"),
		Query(expected.Name, expected.Type, expected.Description, expected.Required),
	)

	assert.Equal(t, 1, len(e.Parameters))
	assert.Equal(t, expected, e.Parameters[0])
}

func TestQueryString(t *testing.T) {
	expected := swag.Parameter{
		In:          "query",
		Name:        "id",
		Description: "the description",
		Required:    false,
		Type:        types.String,
	}

	e := New("get", "/",
		Summary("get thing"),
		QueryS(expected.Name, expected.Description),
	)

	assert.Equal(t, 1, len(e.Parameters))
	assert.Equal(t, expected, e.Parameters[0])
}

func TestFormData(t *testing.T) {
	expected := swag.Parameter{
		In:          "formData",
		Name:        "file",
		Description: "the description",
		Required:    true,
		Type:        types.File,
	}

	e := New("post", "/",
		Summary("upload file"),
		FormData(expected.Name, types.File, expected.Description, expected.Required),
	)

	assert.Equal(t, 1, len(e.Parameters))
	assert.Equal(t, expected, e.Parameters[0])
}

type Model struct {
	String string `json:"s"`
}

func TestBody(t *testing.T) {
	expected := swag.Parameter{
		In:          "body",
		Name:        "body",
		Description: "the description",
		Required:    true,
		Schema: &swag.Schema{
			Ref:       "#/definitions/endpointModel",
			Prototype: reflect.TypeOf(Model{}),
		},
	}

	e := New(
		"get", "/",
		Summary("get thing"),
		Body(Model{}, expected.Description, expected.Required),
	)

	assert.Equal(t, 1, len(e.Parameters))
	assert.Equal(t, expected, e.Parameters[0])
}

func TestResponse(t *testing.T) {
	expected := swag.Response{
		Description: "successful",
		Schema: &swag.Schema{
			Ref:       "#/definitions/endpointModel",
			Prototype: Model{},
		},
	}

	e := New(
		"get", "/",
		Summary("get thing"),
		Response(http.StatusOK, "successful", SchemaResponseOption(Model{})),
	)

	assert.Equal(t, 1, len(e.Responses))
	assert.Equal(t, expected.Description, e.Responses["200"].Description)
	assert.Equal(t, *expected.Schema, *e.Responses["200"].Schema)
}

func TestResponseHeader(t *testing.T) {
	expected := swag.Response{
		Description: "successful",
		Schema:      nil,
		Headers: map[string]swag.Header{
			"X-Rate-Limit": {
				Type:        "integer",
				Format:      "int32",
				Description: "calls per hour allowed by the user",
			},
		},
	}

	e := New(
		"get", "/",
		Summary("get thing"),
		Response(http.StatusOK, "successful",
			HeaderResponseOption("X-Rate-Limit", "integer", "int32", "calls per hour allowed by the user"),
		),
	)

	assert.Equal(t, 1, len(e.Responses))
	assert.Equal(t, expected, e.Responses["200"])
}

func TestSecurityScheme(t *testing.T) {
	api := swag.New(
		option.SecurityScheme("basic", option.BasicSecurity()),
		option.SecurityScheme("apikey", option.APIKeySecurity("Authorization", "header")),
	)
	assert.Len(t, api.SecurityDefinitions, 2)
	assert.Contains(t, api.SecurityDefinitions, "basic")
	assert.Contains(t, api.SecurityDefinitions, "apikey")
	assert.Equal(t, "header", api.SecurityDefinitions["apikey"].In)
}

func TestSecurity(t *testing.T) {
	e := New(
		"get", "/",
		Handler(Echo),
		Security("basic"),
		Security("oauth2", "scope1", "scope2"),
	)
	assert.False(t, e.Security.DisableSecurity)
	assert.Len(t, e.Security.Requirements, 2)
	assert.Contains(t, e.Security.Requirements[0], "basic")
	assert.Contains(t, e.Security.Requirements[1], "oauth2")
	assert.Len(t, e.Security.Requirements[1]["oauth2"], 2)
}

func TestNoSecurity(t *testing.T) {
	e := New(
		"get", "/",
		Handler(Echo),
		NoSecurity(),
	)
	assert.True(t, e.Security.DisableSecurity)
}
