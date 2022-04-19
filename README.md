# swag

![LICENSE](https://img.shields.io/github/license/zc2638/swag.svg?style=flat-square&color=blue)
[![GoDoc](https://godoc.org/github.com/zc2638/swag?status.svg)](https://godoc.org/github.com/zc2638/swag)
[![Go Report Card](https://goreportcard.com/badge/github.com/zc2638/swag)](https://goreportcard.com/report/github.com/zc2638/swag)


```swag``` is a lightweight library to generate swagger json for Go projects.  
 
No code generation, no framework constraints, just a simple swagger definition.

```swag``` is heavily geared towards generating REST/JSON apis.

## Dependency

Golang 1.16+

## Installation

```shell
go get -u github.com/zc2638/swag
```

## Concepts

```swag``` uses functional options to generate both the swagger endpoints and the swagger definition.  Where possible
```swag``` attempts to use reasonable defaults that may be overridden by the user.

### Endpoints

```swag``` provides a separate package, ```endpoint```, to generate swagger endpoints.  These endpoints can be passed
to the swagger definition generate via ```swag.Endpoints(...)```

In this simple example, we generate an endpoint to retrieve all pets.  The only required fields for an endpoint
are the method, path, and the summary.  

```go
allPets := endpoint.New("get", "/pet", "Return all the pets") 
```

However, it'll probably be useful if you include definitions of what ```GET /pet``` returns:

```go
allPets := endpoint.New("get", "/pet", "Return all the pets",
  endpoint.Response(http.StatusOk, Pet{}, "Successful operation"),
  endpoint.Response(http.StatusInternalServerError, Error{}, "Oops ... something went wrong"),
) 
```

Refer to the [godoc](https://godoc.org/github.com/zc2638/swag/endpoint) for a list of all the endpoint options

### Walk

As a convenience to users, ```*swag.API``` implements a ```Walk``` method to simplify traversal of all the endpoints.
See the complete example below for how ```Walk``` can be used to bind endpoints to the router.

```go
api := swag.New(
    swag.TitleOption("Swagger Petstore"),
    swag.EndpointsOption(post, get),
)

// iterate over each endpoint, if we've defined a handler, we can use it to bind to the router.  We're using ```gin``
// in this example, but any web framework will do.
// 
api.Walk(func(path string, endpoint *swag.Endpoint) {
    h := endpoint.Handler.(func(c *gin.Context))
    path = swag.ColonPath(path)
    router.Handle(endpoint.Method, path, h)
})
```

## Example
### definition
```go
package main

import (
	"fmt"
	"io"
	"net/http"
	"path"

	"github.com/zc2638/swag"
	"github.com/zc2638/swag/endpoint"
)

// Category example from the swagger pet store
type Category struct {
	ID     int64  `json:"category"`
	Name   string `json:"name" enum:"dog,cat" required:""`
	Exists *bool  `json:"exists" required:""`
}

// Pet example from the swagger pet store
type Pet struct {
	ID        int64     `json:"id"`
	Category  *Category `json:"category" desc:"分类"`
	Name      string    `json:"name" required:"" example:"张三" desc:"名称"`
	PhotoUrls []string  `json:"photoUrls"`
	Tags      []string  `json:"tags" desc:"标签"`
}

func handle(w http.ResponseWriter, r *http.Request) {
	_, _ = io.WriteString(w, fmt.Sprintf("[%s] Hello World!", r.Method))
}

func main() {
	post := endpoint.New("post", "/pet", endpoint.Summary("Add a new pet to the store"),
		endpoint.Handler(handle),
		endpoint.Description("Additional information on adding a pet to the store"),
		endpoint.Body(Pet{}, "Pet object that needs to be added to the store", true),
		endpoint.Response(http.StatusOK, "Successfully added pet", endpoint.Schema(Pet{})),
		endpoint.Security("petstore_auth", "read:pets", "write:pets"),
	)
	get := endpoint.New("get", "/pet/{petId}", endpoint.Summary("Find pet by ID"),
		endpoint.Handler(handle),
		endpoint.Path("petId", "integer", "ID of pet to return", true),
		endpoint.Response(http.StatusOK, "successful operation", endpoint.Schema(Pet{})),
		endpoint.Security("petstore_auth", "read:pets"),
	)
	test := endpoint.New("put", "/pet/{petId}",
		endpoint.Handler(handle),
		endpoint.Path("petId", "integer", "ID of pet to return", true),
		endpoint.Response(http.StatusOK, "successful operation", endpoint.Schema(struct {
			ID   string `json:"id"`
			Name string `json:"name"`
		}{})),
		endpoint.Security("petstore_auth", "read:pets"),
	)

	securityScheme := swag.OAuth2Security("accessCode", "http://example.com/oauth/authorize", "http://example.com/oauth/token")
	securityScheme.Scopes["write:pets"] = "modify pets in your account"
	securityScheme.Scopes["read:pets"] = "read your pets"
	api := swag.New(
		swag.TitleOption("Example API Doc"),
		swag.SecurityOption("petstore_auth", "read:pets"),
		swag.SecuritySchemeOption("petstore_auth", securityScheme),
		swag.EndpointsOption(post, get),
	)
	api.AddEndpoint(test)

	...
}

```

### builtin
```go
func main() {
	...

	for p, endpoints := range api.Paths {
		http.DefaultServeMux.Handle(path.Join(api.BasePath, p), endpoints)
	}
	http.DefaultServeMux.Handle("/swagger/json", api.Handler())
	patterns := swag.UIPatterns("/swagger/ui")
	for _, pattern := range patterns {
		http.DefaultServeMux.Handle(pattern, swag.UIHandler("/swagger/ui", "/swagger/json", true))
	}

	http.ListenAndServe(":8080", nil)
}

```

### gin
```go
func main() {
    ...
	
    router := gin.New()
    api.Walk(func(path string, endpoint *swag.Endpoint) {
        h := endpoint.Handler.(func(c *gin.Context))
        path = swag.ColonPath(path)
        
        router.Handle(endpoint.Method, path, h)
    })
    
    // Register Swagger JSON route
    router.GET("/swagger/json", gin.WrapH(api.Handler()))
    
    // Register Swagger UI route
    // To take effect, the swagger json route must be registered
    router.GET("/swagger/ui/*any",
        gin.WrapH(swag.UIHandler("/swagger/ui", "/swagger/json", true)))
    
    http.ListenAndServe(":8080", router)
}
```

### mux
```go
func main() {
	...

	router := mux.NewRouter()
	api.Walk(func(path string, endpoint *swag.Endpoint) {
		h := endpoint.Handler.(http.HandlerFunc)
		path = swag.ColonPath(path)

		router.Path(path).Methods(endpoint.Method).Handler(h)
	})

	router.Path("/swagger/json").Methods("GET").Handler(api.Handler())
	router.PathPrefix("/swagger/ui").Handler(
		swag.UIHandler("/swagger/ui", "/swagger/json", true))

	http.ListenAndServe(":8080", router)
}
```

### echo
```go
func main() {
	...

	router := echo.New()
	api.Walk(func(path string, endpoint *swag.Endpoint) {
		h := endpoint.Handler.(func(c echo.Context) error)
		path = swag.ColonPath(path)

		switch strings.ToLower(endpoint.Method) {
		case "get":
			router.GET(path, h)
		case "head":
			router.HEAD(path, h)
		case "options":
			router.OPTIONS(path, h)
		case "delete":
			router.DELETE(path, h)
		case "put":
			router.PUT(path, h)
		case "post":
			router.POST(path, h)
		case "trace":
			router.TRACE(path, h)
		case "patch":
			router.PATCH(path, h)
		case "connect":
			router.CONNECT(path, h)
		}
	})

	router.GET("/swagger/json", echo.WrapHandler(api.Handler()))
	router.GET("/swagger/ui/*",
		echo.WrapHandler(swag.UIHandler("/swagger/ui", "/swagger/json", true)))

	http.ListenAndServe(":8080", router)
}
```

### httprouter
```go
func main() {
	...

	router := httprouter.New()
	api.Walk(func(path string, endpoint *swag.Endpoint) {
		h := endpoint.Handler.(func(w http.ResponseWriter, req *http.Request, param httprouter.Params))
		path = swag.ColonPath(path)
		router.Handle(endpoint.Method, path, h)
	})

	router.Handler(http.MethodGet, "/swagger/json", api.Handler())
	router.Handler(http.MethodGet, "/swagger/ui/*any",
		swag.UIHandler("/swagger/ui", "/swagger/json", true))

	http.ListenAndServe(":8080", router)
}
```

## Default Swagger UI Server

```go
func main() {
    http.DefaultServeMux.Handle("/swagger/json", api.Handler())
    patterns := swag.UIPatterns("/swagger/ui")
    handle := swag.UIHandler("/swagger/ui", "", false)
    for _, pattern := range patterns {
        http.DefaultServeMux.Handle(pattern, handle)
    }
}
```
so you can visit for config: `http://localhost:8080/swagger/json`  
so you can visit for ui: `http://localhost:8080/swagger/ui`
