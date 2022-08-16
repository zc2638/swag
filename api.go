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

package swag

import (
	"bytes"
	"fmt"
	"io"
	"net/http"
	"path"
	"path/filepath"
	"reflect"
	"strings"

	"github.com/99nil/gopkg/ctr"

	"github.com/zc2638/swag/asserts"
)

// Object represents the object entity from the swagger definition
type Object struct {
	IsArray     bool                `json:"-"`
	GoType      reflect.Type        `json:"-"`
	Name        string              `json:"-"`
	Type        string              `json:"type"`
	Description string              `json:"description,omitempty"`
	Format      string              `json:"format,omitempty"`
	Required    []string            `json:"required,omitempty"`
	Properties  map[string]Property `json:"properties,omitempty"`
}

// Property represents the property entity from the swagger definition
type Property struct {
	GoType      reflect.Type `json:"-"`
	Type        string       `json:"type,omitempty"`
	Description string       `json:"description,omitempty"`
	Enum        []string     `json:"enum,omitempty"`
	Format      string       `json:"format,omitempty"`
	Ref         string       `json:"$ref,omitempty"`
	Example     string       `json:"example,omitempty"`
	Items       *Items       `json:"items,omitempty"`
}

// Contact represents the contact entity from the swagger definition; used by Info
type Contact struct {
	Email string `json:"email,omitempty"`
}

// License represents the license entity from the swagger definition; used by Info
type License struct {
	Name string `json:"name,omitempty"`
	URL  string `json:"url,omitempty"`
}

// Info represents the info entity from the swagger definition
type Info struct {
	Description    string   `json:"description,omitempty"`
	Version        string   `json:"version,omitempty"`
	TermsOfService string   `json:"termsOfService,omitempty"`
	Title          string   `json:"title,omitempty"`
	Contact        *Contact `json:"contact,omitempty"`
	License        License  `json:"license"`
}

// SecurityScheme represents a security scheme from the swagger definition.
type SecurityScheme struct {
	Type             string            `json:"type"`
	Description      string            `json:"description,omitempty"`
	Name             string            `json:"name,omitempty"`
	In               string            `json:"in,omitempty"`
	Flow             string            `json:"flow,omitempty"`
	AuthorizationURL string            `json:"authorizationUrl,omitempty"`
	TokenURL         string            `json:"tokenUrl,omitempty"`
	Scopes           map[string]string `json:"scopes,omitempty"`
}

// Endpoints represents all the swagger endpoints associated with a particular path
type Endpoints struct {
	Delete  *Endpoint `json:"delete,omitempty"`
	Head    *Endpoint `json:"head,omitempty"`
	Get     *Endpoint `json:"get,omitempty"`
	Options *Endpoint `json:"options,omitempty"`
	Post    *Endpoint `json:"post,omitempty"`
	Put     *Endpoint `json:"put,omitempty"`
	Patch   *Endpoint `json:"patch,omitempty"`
	Trace   *Endpoint `json:"trace,omitempty"`
	Connect *Endpoint `json:"connect,omitempty"`
}

// ServeHTTP allows endpoints to serve itself using the builtin http mux
func (e *Endpoints) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	var endpoint *Endpoint

	switch req.Method {
	case http.MethodDelete:
		endpoint = e.Delete
	case http.MethodHead:
		endpoint = e.Head
	case http.MethodGet:
		endpoint = e.Get
	case http.MethodOptions:
		endpoint = e.Options
	case http.MethodPost:
		endpoint = e.Post
	case http.MethodPut:
		endpoint = e.Put
	case http.MethodPatch:
		endpoint = e.Patch
	case http.MethodTrace:
		endpoint = e.Trace
	case http.MethodConnect:
		endpoint = e.Connect
	}

	if endpoint == nil || endpoint.Handler == nil {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	switch v := endpoint.Handler.(type) {
	case func(w http.ResponseWriter, req *http.Request):
		v(w, req)
	case http.HandlerFunc:
		v(w, req)
	case http.Handler:
		v.ServeHTTP(w, req)
	default:
		w.WriteHeader(http.StatusInternalServerError)
		_, _ = io.WriteString(w, "Handler is not a standard http handler")
	}
}

// Walk calls the specified function for each method defined within the Endpoints
func (e *Endpoints) Walk(fn func(endpoint *Endpoint)) {
	if e.Delete != nil {
		fn(e.Delete)
	}
	if e.Head != nil {
		fn(e.Head)
	}
	if e.Get != nil {
		fn(e.Get)
	}
	if e.Options != nil {
		fn(e.Options)
	}
	if e.Post != nil {
		fn(e.Post)
	}
	if e.Put != nil {
		fn(e.Put)
	}
	if e.Patch != nil {
		fn(e.Patch)
	}
	if e.Trace != nil {
		fn(e.Trace)
	}
	if e.Connect != nil {
		fn(e.Connect)
	}
}

// API provides the top level encapsulation for the swagger definition
type API struct {
	Swagger             string                    `json:"swagger,omitempty"`
	Info                Info                      `json:"info"`
	BasePath            string                    `json:"basePath,omitempty"`
	Schemes             []string                  `json:"schemes,omitempty"`
	Paths               map[string]*Endpoints     `json:"paths,omitempty"`
	Definitions         map[string]Object         `json:"definitions,omitempty"`
	Tags                []Tag                     `json:"tags,omitempty"`
	Host                string                    `json:"host,omitempty"`
	SecurityDefinitions map[string]SecurityScheme `json:"securityDefinitions,omitempty"`
	Security            *SecurityRequirement      `json:"security,omitempty"`
	tags                []Tag
}

func (a *API) Clone() *API {
	return &API{
		Swagger:             a.Swagger,
		Info:                a.Info,
		BasePath:            a.BasePath,
		Schemes:             a.Schemes,
		Paths:               a.Paths,
		Definitions:         a.Definitions,
		Tags:                a.Tags,
		Host:                a.Host,
		SecurityDefinitions: a.SecurityDefinitions,
		Security:            a.Security,
	}
}

func (a *API) addPath(e *Endpoint) {
	if a.Paths == nil {
		a.Paths = map[string]*Endpoints{}
	}

	v, ok := a.Paths[e.Path]
	if !ok {
		v = &Endpoints{}
		a.Paths[e.Path] = v
	}

	switch strings.ToUpper(e.Method) {
	case http.MethodDelete:
		v.Delete = e
	case http.MethodGet:
		v.Get = e
	case http.MethodHead:
		v.Head = e
	case http.MethodOptions:
		v.Options = e
	case http.MethodPost:
		v.Post = e
	case http.MethodPut:
		v.Put = e
	case http.MethodPatch:
		v.Patch = e
	case http.MethodTrace:
		v.Trace = e
	case http.MethodConnect:
		v.Connect = e
	default:
		panic(fmt.Errorf("invalid method, %v", e.Method))
	}
}

func (a *API) addDefinition(e *Endpoint) {
	if a.Definitions == nil {
		a.Definitions = map[string]Object{}
	}

	if e.Parameters != nil {
		for _, p := range e.Parameters {
			if p.Schema != nil {
				def := define(p.Schema.Prototype)
				for k, v := range def {
					if _, ok := a.Definitions[k]; !ok {
						a.Definitions[k] = v
					}
				}
			}
		}
	}

	if e.Responses != nil {
		for _, response := range e.Responses {
			if response.Schema != nil {
				def := define(response.Schema.Prototype)
				for k, v := range def {
					if _, ok := a.Definitions[k]; !ok {
						a.Definitions[k] = v
					}
				}
			}
		}
	}
}

func (a *API) WithTags(tags ...Tag) *API {
	for _, v := range tags {
		exists := false
		for _, tag := range a.Tags {
			if tag.Name == v.Name {
				exists = true
				break
			}
		}
		if exists {
			continue
		}
		a.Tags = append(a.Tags, v)
	}
	a.tags = tags
	return a
}

// AddEndpoint adds the specified endpoint to the API definition; to generate an endpoint use ```endpoint.New```
func (a *API) AddEndpoint(es ...*Endpoint) {
	tags := make([]string, 0, len(a.tags))
	for _, tag := range a.tags {
		tags = append(tags, tag.Name)
	}
	for _, e := range es {
		e.Tags = append(e.Tags, tags...)
		a.addPath(e)
		a.addDefinition(e)
	}
	a.tags = nil
}

// AddOptions adds some options
func (a *API) AddOptions(options ...Option) {
	for _, option := range options {
		option(a)
	}
}

// AddEndpointFunc adds some options
//
// Deprecated: please use the new AddOptions method
func (a *API) AddEndpointFunc(fs ...func(*API)) {
	for _, f := range fs {
		f(a)
	}
	a.tags = nil
}

func (a *API) AddTag(name, description string) {
	a.Tags = append(a.Tags, Tag{
		Name:        name,
		Description: description,
	})
}

// Handler is a factory method that generates an http.HandlerFunc; if enableCors is true, then the handler will generate
// cors headers
func (a *API) Handler() http.HandlerFunc {
	return func(w http.ResponseWriter, req *http.Request) {
		// customize the swagger header based on host
		scheme := ""
		if req.TLS != nil {
			scheme = "https"
		}
		if v := req.Header.Get("X-Forwarded-Proto"); v != "" {
			scheme = v
		}
		if scheme == "" {
			scheme = req.URL.Scheme
		}
		if scheme == "" {
			scheme = "http"
		}
		doc := a.Clone()
		doc.Host = req.Host
		doc.Schemes = []string{scheme}
		ctr.OK(w, doc)
	}
}

// Walk invoke the callback for each endpoints defined in the swagger doc
func (a *API) Walk(callback func(path string, endpoints *Endpoint)) {
	for rawPath, endpoints := range a.Paths {
		u := path.Join(a.BasePath, rawPath)
		endpoints.Walk(func(endpoint *Endpoint) {
			callback(u, endpoint)
		})
	}
}

// UIPatterns returns a list of all the paths needed based on the path prefix
func UIPatterns(prefix string) []string {
	files, err := asserts.Dist.ReadDir(asserts.DistDir)
	if err != nil {
		return nil
	}
	patterns := make([]string, 0, len(files)+1)
	patterns = append(patterns, path.Join(prefix)+"/")
	for _, f := range files {
		patterns = append(patterns, path.Join(prefix, f.Name()))
	}
	return patterns
}

// UIHandler returns a http.Handler by the specify path prefix and the full path
func UIHandler(prefix, uri string, autoDomain bool) http.Handler {
	return http.StripPrefix(prefix, http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "" {
			url := strings.TrimSuffix(prefix, "/") + "/"
			http.Redirect(w, r, url, http.StatusFound)
			return
		}
		if r.URL.Path == "/" || r.URL.Path == "index.html" {
			fullName := filepath.Join(asserts.DistDir, "index.html")
			fileData, err := asserts.Dist.ReadFile(fullName)
			if err != nil {
				w.WriteHeader(http.StatusInternalServerError)
				_, _ = w.Write([]byte("index.html read exception"))
				return
			}
			if uri == "" {
				_, _ = w.Write(fileData)
				return
			}

			// Prevent uri assignment from causing final uri exception.
			currentURI := uri
			if autoDomain {
				scheme := ""
				if r.TLS != nil {
					scheme = "https"
				}
				if v := r.Header.Get("X-Forwarded-Proto"); v != "" {
					scheme = v
				}
				if scheme == "" {
					scheme = r.URL.Scheme
				}
				if scheme == "" {
					scheme = "http"
				}
				currentURI = scheme + "://" + path.Join(r.Host, currentURI)
			}

			fileData = bytes.ReplaceAll(fileData, []byte(asserts.URL), []byte(currentURI))
			_, _ = w.Write(fileData)
			return
		}
		http.FileServer(DirFS(asserts.DistDir, asserts.Dist)).ServeHTTP(w, r)
	}))
}
