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
//
package swagger

import (
	"bytes"
	"fmt"
	"github.com/pkgms/go/ctr"
	"github.com/zc2638/swag/asserts"
	"io"
	"io/fs"
	"net/http"
	"path"
	"path/filepath"
	"reflect"
	"strings"
)

// Object represents the object entity from the swagger definition
type Object struct {
	IsArray    bool                `json:"-"`
	GoType     reflect.Type        `json:"-"`
	Name       string              `json:"-"`
	Type       string              `json:"type"`
	Format     string              `json:"format,omitempty"`
	Required   []string            `json:"required,omitempty"`
	Properties map[string]Property `json:"properties,omitempty"`
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

// SecuritySchemeOption provides additional customizations to the SecurityScheme.
type SecuritySchemeOption func(securityScheme *SecurityScheme)

// SecuritySchemeDescription sets the security scheme's description.
func SecuritySchemeDescription(description string) SecuritySchemeOption {
	return func(securityScheme *SecurityScheme) {
		securityScheme.Description = description
	}
}

// BasicSecurity defines a security scheme for HTTP Basic authentication.
func BasicSecurity() SecuritySchemeOption {
	return func(securityScheme *SecurityScheme) {
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

	return func(securityScheme *SecurityScheme) {
		securityScheme.Type = "apiKey"
		securityScheme.Name = name
		securityScheme.In = in
	}
}

// OAuth2Scope adds a new scope to the security scheme.
func OAuth2Scope(scope, description string) SecuritySchemeOption {
	return func(securityScheme *SecurityScheme) {
		if securityScheme.Scopes == nil {
			securityScheme.Scopes = map[string]string{}
		}
		securityScheme.Scopes[scope] = description
	}
}

// OAuth2Security defines a security scheme for OAuth2 authentication. Flow can
// be one of implicit, password, application, or accessCode.
func OAuth2Security(flow, authorizationURL, tokenURL string) SecuritySchemeOption {
	return func(securityScheme *SecurityScheme) {
		securityScheme.Type = "oauth2"
		securityScheme.Flow = flow
		securityScheme.AuthorizationURL = authorizationURL
		securityScheme.TokenURL = tokenURL
		if securityScheme.Scopes == nil {
			securityScheme.Scopes = map[string]string{}
		}
	}
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
		io.WriteString(w, "Handler is not a standard http handler")
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
}

func (a *API) clone() *API {
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

// AddEndpoint adds the specified endpoint to the API definition; to generate an endpoint use ```endpoint.New```
func (a *API) AddEndpoint(es ...*Endpoint) {
	for _, e := range es {
		a.addPath(e)
		a.addDefinition(e)
	}
}

func (a *API) AddEndpointFunc(fs ...func(*API)) {
	for _, f := range fs {
		f(a)
	}
}

func (a *API) AddTag(name, description string) {
	a.Tags = append(a.Tags, Tag{
		Name:        name,
		Description: description,
	})
}

// Handler is a factory method that generates an http.HandlerFunc; if enableCors is true, then the handler will generate
// cors headers
func (a *API) Handler(enableCors bool) http.HandlerFunc {
	return func(w http.ResponseWriter, req *http.Request) {
		if enableCors {
			w.Header().Set("Access-Control-Allow-Headers", "Content-Type, api_key, Authorization")
			w.Header().Set("Access-Control-Allow-Methods", "GET, POST, DELETE, PUT, PATCH")
			w.Header().Set("Access-Control-Allow-Origin", "*")
		}

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
		doc := a.clone()
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

type RouteInterface interface {
	Handle(pattern string, handler http.Handler)
}

func (a *API) registerMux(router RouteInterface, url string, autoDomain bool) {
	files, err := asserts.Dist.ReadDir(asserts.DistDir)
	if err != nil {
		return
	}
	handler := http.StripPrefix("/swagger-ui", http.FileServer(DirFS(asserts.DistDir, asserts.Dist)))
	for _, file := range files {
		filename := file.Name()
		pattern := path.Join("/swagger-ui", filename)
		if filename == "index.html" {
			fullName := filepath.Join(asserts.DistDir, filename)
			fileData, err := asserts.Dist.ReadFile(fullName)
			if err != nil {
				return
			}
			if url == "" {
				url = asserts.URL
			}
			indexHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
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
					url = scheme + "://" + path.Join(r.Host, url)
				}
				fileData = bytes.ReplaceAll(fileData, []byte(asserts.URL), []byte(url))
				w.Write(fileData)
			})
			router.Handle(pattern, indexHandler)
			router.Handle("/swagger-ui", http.RedirectHandler("/swagger-ui/index.html", http.StatusFound))
			router.Handle("/swagger-ui/", http.RedirectHandler("/swagger-ui/index.html", http.StatusFound))
			continue
		}
		router.Handle(pattern, handler)
	}
}

func (a *API) RegisterMux(router RouteInterface, url string) {
	a.registerMux(router, url, false)
}

func (a *API) RegisterMuxWithData(router RouteInterface, enableCors bool) {
	for p, endpoints := range a.Paths {
		router.Handle(p, endpoints)
	}
	const url = "/swagger-ui/json"
	router.Handle(url, a.Handler(enableCors))
	a.registerMux(router, url, true)
}

func DirFS(dir string, fsys fs.FS) http.FileSystem {
	return dirFS{
		dir: dir,
		fs:  http.FS(fsys),
	}
}

type dirFS struct {
	dir string
	fs  http.FileSystem
}

func (f dirFS) Open(name string) (http.File, error) {
	return f.fs.Open(filepath.Join(f.dir, name))
}
