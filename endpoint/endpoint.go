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
	"net/http"
	"reflect"
	"strconv"
	"strings"

	"github.com/99nil/gopkg/sets"

	"github.com/zc2638/swag/types"

	"github.com/zc2638/swag"
)

// New constructs a new swagger endpoint using the fields and functional options provided
func New(method, path string, options ...Option) *swag.Endpoint {
	e := &swag.Endpoint{
		Method:   strings.ToUpper(method),
		Path:     path,
		Produces: []string{"application/json"},
		Consumes: []string{"application/json"},
	}
	e.BuildOperationID()

	for _, opt := range options {
		opt(e)
	}
	return e
}

// Option represents a functional option to customize the swagger endpoint
type Option func(e *swag.Endpoint)

// Handler allows an instance of the web handler to be associated with the endpoint.  This can be especially useful when
// using swag to bind the endpoints to the web router.  See the examples package for how the Handler can be used in
// conjunction with Walk to simplify binding endpoints to a router
func Handler(handler interface{}) Option {
	return func(e *swag.Endpoint) {
		if v, ok := handler.(func(w http.ResponseWriter, r *http.Request)); ok {
			handler = http.HandlerFunc(v)
		}
		e.Handler = handler
	}
}

// Summary sets the endpoint's summary
func Summary(v string) Option {
	return func(e *swag.Endpoint) {
		e.Summary = v
	}
}

// Description sets the endpoint's description
func Description(v string) Option {
	return func(e *swag.Endpoint) {
		e.Description = v
	}
}

// OperationID sets the endpoint's operationId
func OperationID(v string) Option {
	return func(e *swag.Endpoint) {
		e.OperationID = v
	}
}

// Produces sets the endpoint's produces; by default this will be set to application/json
func Produces(v ...string) Option {
	return func(e *swag.Endpoint) {
		e.Produces = v
	}
}

// Consumes sets the endpoint's produces; by default this will be set to application/json
func Consumes(v ...string) Option {
	return func(e *swag.Endpoint) {
		e.Consumes = v
	}
}

func parameter(p swag.Parameter) Option {
	return func(e *swag.Endpoint) {
		if e.Parameters == nil {
			e.Parameters = make([]swag.Parameter, 0)
		}
		e.Parameters = append(e.Parameters, p)
	}
}

// Path defines a path parameter for the endpoint;
// name, typ, description and required correspond to the matching swagger fields
func Path(name string, typ types.ParameterType, description string, required bool) Option {
	return PathDefault(name, typ, description, "", required)
}

// PathString defines a path parameter for the endpoint;
// name and description correspond to the matching swagger fields,
// type defaults to string,
// required defaults to true.
func PathString(name, description string) Option {
	return PathDefault(name, types.String, description, "", true)
}

// PathDefault defines a path parameter for the endpoint;
// name, typ, description, defVal and required correspond to the matching swagger fields
func PathDefault(name string, typ types.ParameterType, description, defVal string, required bool) Option {
	p := swag.Parameter{
		Name:        name,
		In:          "path",
		Type:        typ,
		Description: description,
		Required:    required,
		Default:     defVal,
	}
	return parameter(p)
}

// Query defines a query parameter for the endpoint;
// name, typ, description and required correspond to the matching swagger fields
func Query(name string, typ types.ParameterType, description string, required bool) Option {
	return QueryDefault(name, typ, description, "", required)
}

// QueryString defines a query parameter for the endpoint;
// name and description correspond to the matching swagger fields,
// type defaults to string,
// required defaults to false.
func QueryString(name, description string) Option {
	return QueryDefault(name, types.String, description, "", false)
}

// QueryDefault defines a query parameter for the endpoint;
// name, typ, description, defVal and required correspond to the matching swagger fields
func QueryDefault(name string, typ types.ParameterType, description, defVal string, required bool) Option {
	p := swag.Parameter{
		Name:        name,
		In:          "query",
		Type:        typ,
		Description: description,
		Required:    required,
		Default:     defVal,
	}
	return parameter(p)
}

// FormData defines a form-data parameter for the endpoint;
// name, typ, description and required correspond to the matching swagger fields
func FormData(name string, typ types.ParameterType, description string, required bool) Option {
	p := swag.Parameter{
		In:          "formData",
		Type:        typ,
		Name:        name,
		Description: description,
		Required:    required,
	}
	return func(e *swag.Endpoint) {
		parameter(p)(e)

		s := sets.NewString(e.Consumes...)
		s.Add("multipart/form-data")
		e.Consumes = s.List()
	}
}

// Body defines a body parameter for the swagger endpoint as would commonly be used for the POST, PUT, and PATCH methods
// prototype should be a struct or a pointer to struct that swag can use to reflect upon the return type
func Body(prototype interface{}, description string, required bool) Option {
	return BodyType(reflect.TypeOf(prototype), description, required)
}

// BodyType defines a body parameter for the swagger endpoint as would commonly be used for the POST, PUT, and PATCH methods
// prototype should be a struct or a pointer to struct that swag can use to reflect upon the return type
// t represents the Type of the body
func BodyType(t reflect.Type, description string, required bool) Option {
	p := swag.Parameter{
		In:          "body",
		Name:        "body",
		Description: description,
		Schema:      swag.MakeSchema(t),
		Required:    required,
	}
	return parameter(p)
}

// Tags allows one or more tags to be associated with the endpoint
func Tags(tags ...string) Option {
	return func(e *swag.Endpoint) {
		if e.Tags == nil {
			e.Tags = make([]string, 0)
		}
		e.Tags = append(e.Tags, tags...)
	}
}

// Security allows a security scheme to be associated with the endpoint.
func Security(scheme string, scopes ...string) Option {
	return func(e *swag.Endpoint) {
		if e.Security == nil {
			e.Security = &swag.SecurityRequirement{}
		}

		if e.Security.Requirements == nil {
			e.Security.Requirements = []map[string][]string{}
		}

		e.Security.Requirements = append(e.Security.Requirements, map[string][]string{scheme: scopes})
	}
}

// NoSecurity explicitly sets the endpoint to have no security requirements.
func NoSecurity() Option {
	return func(e *swag.Endpoint) {
		e.Security = &swag.SecurityRequirement{DisableSecurity: true}
	}
}

// ResponseOption allows for additional configurations on responses like header information
type ResponseOption func(response *swag.Response)

// Schema adds schema definitions to swagger responses
func Schema(schema interface{}) ResponseOption {
	return func(response *swag.Response) {
		response.Schema = swag.MakeSchema(schema)
	}
}

// Header adds header definitions to swagger responses
func Header(name, typ, format, description string) ResponseOption {
	return func(response *swag.Response) {
		if response.Headers == nil {
			response.Headers = map[string]swag.Header{}
		}
		response.Headers[name] = swag.Header{
			Type:        typ,
			Format:      format,
			Description: description,
		}
	}
}

// ResponseType sets the endpoint response for the specified code; may be used multiple times with different status codes
// t represents the Type of the response
func ResponseType(code int, description string, opts ...ResponseOption) Option {
	return func(e *swag.Endpoint) {
		if e.Responses == nil {
			e.Responses = make(map[string]swag.Response)
		}
		r := swag.Response{
			Description: description,
		}
		for _, opt := range opts {
			opt(&r)
		}
		e.Responses[strconv.Itoa(code)] = r
	}
}

// Response sets the endpoint response for the specified code; may be used multiple times with different status codes
func Response(code int, description string, opts ...ResponseOption) Option {
	return ResponseType(code, description, opts...)
}

func ResponseSuccess(opts ...ResponseOption) Option {
	return Response(http.StatusOK, "success", opts...)
}

func Deprecated() Option {
	return func(e *swag.Endpoint) {
		e.Deprecated = true
	}
}
