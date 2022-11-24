# swag

![LICENSE](https://img.shields.io/github/license/zc2638/swag.svg?style=flat-square&color=blue)
[![GoDoc](https://godoc.org/github.com/zc2638/swag?status.svg)](https://godoc.org/github.com/zc2638/swag)
[![Go Report Card](https://goreportcard.com/badge/github.com/zc2638/swag?style=flat-square)](https://goreportcard.com/report/github.com/zc2638/swag)

[English](./README.md) | 简体中文

```swag``` 是一个轻量级的库，用于为 Golang 项目生成 `Swagger JSON`。

```swag``` 主要用于生成 REST/JSON API接口。

没有代码生成，没有框架约束，只是一个简单的 swagger 定义。

## 依赖

Golang 1.16+

## 安装

```shell
go get -u github.com/zc2638/swag
```

**Tip:** 从 `v1.2.0` 开始，低版本不再兼容。为了兼容大部分的web框架，整体架构做了很大的改动。

## 默认 Swagger UI 服务器

```go
func main() {
    handle := swag.UIHandler("/swagger/ui", "", false)
    patterns := swag.UIPatterns("/swagger/ui")
    for _, pattern := range patterns {
        http.DefaultServeMux.Handle(pattern, handle)
    }

    log.Fatal(http.ListenAndServe(":8080", nil))
}
```
可以通过此地址访问 UI: `http://localhost:8080/swagger/ui`

## Examples

### 定义

```go
package main

import (
	"fmt"
	"io"
	"net/http"

	"github.com/zc2638/swag"
	"github.com/zc2638/swag/endpoint"
	"github.com/zc2638/swag/option"
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
	api := swag.New(
		option.Title("Example API Doc"),
		option.Security("petstore_auth", "read:pets"),
		option.SecurityScheme("petstore_auth",
			option.OAuth2Security("accessCode", "http://example.com/oauth/authorize", "http://example.com/oauth/token"),
			option.OAuth2Scope("write:pets", "modify pets in your account"),
			option.OAuth2Scope("read:pets", "read your pets"),
		),
	)
	api.AddEndpoint(
		endpoint.New(
			http.MethodPost, "/pet",
			endpoint.Handler(handle),
			endpoint.Summary("Add a new pet to the store"),
			endpoint.Description("Additional information on adding a pet to the store"),
			endpoint.Body(Pet{}, "Pet object that needs to be added to the store", true),
			endpoint.Response(http.StatusOK, "Successfully added pet", endpoint.Schema(Pet{})),
			endpoint.Security("petstore_auth", "read:pets", "write:pets"),
		),
		endpoint.New(
			http.MethodGet, "/pet/{petId}",
			endpoint.Handler(handle),
			endpoint.Summary("Find pet by ID"),
			endpoint.Path("petId", "integer", "ID of pet to return", true),
			endpoint.Response(http.StatusOK, "successful operation", endpoint.Schema(Pet{})),
			endpoint.Security("petstore_auth", "read:pets"),
		),
		endpoint.New(
			http.MethodPut, "/pet/{petId}",
			endpoint.Handler(handle),
			endpoint.Path("petId", "integer", "ID of pet to return", true),
			endpoint.Security("petstore_auth", "read:pets"),
			endpoint.ResponseSuccess(endpoint.Schema(struct {
				ID   string `json:"id"`
				Name string `json:"name"`
			}{})),
		),
	)

	...
}

```

### built-in

```go
func main() {
    ...
    // Note: Built-in routes cannot automatically resolve path parameters.
    for p, endpoints := range api.Paths {
        http.DefaultServeMux.Handle(path.Join(api.BasePath, p), endpoints)
    }
    http.DefaultServeMux.Handle("/swagger/json", api.Handler())
    patterns := swag.UIPatterns("/swagger/ui")
    for _, pattern := range patterns {
        http.DefaultServeMux.Handle(pattern, swag.UIHandler("/swagger/ui", "/swagger/json", true))
    }

    log.Fatal(http.ListenAndServe(":8080", nil))
}

```

### gin

```go
func main() {
    ...

    router := gin.New()
    api.Walk(func (path string, e *swag.Endpoint) {
        h := e.Handler.(http.Handler)
        path = swag.ColonPath(path)

        router.Handle(e.Method, path, gin.WrapH(h))
    })
    
    // Register Swagger JSON route
    router.GET("/swagger/json", gin.WrapH(api.Handler()))
    
    // Register Swagger UI route
    // To take effect, the swagger json route must be registered
    router.GET("/swagger/ui/*any", gin.WrapH(swag.UIHandler("/swagger/ui", "/swagger/json", true)))

    log.Fatal(http.ListenAndServe(":8080", router))
}
```

### chi

```go
func main() {
    ...
    
    router := chi.NewRouter()
    api.Walk(func (path string, e *swag.Endpoint) {
        router.Method(e.Method, path, e.Handler.(http.Handler))
    })
    router.Handle("/swagger/json", api.Handler())
    router.Mount("/swagger/ui", swag.UIHandler("/swagger/ui", "/swagger/json", true))

    log.Fatal(http.ListenAndServe(":8080", router))
}
```

### mux

```go
func main() {
    ...
    
    router := mux.NewRouter()
    api.Walk(func (path string, e *swag.Endpoint) {
        h := e.Handler.(http.HandlerFunc)
        router.Path(path).Methods(e.Method).Handler(h)
    })
    
    router.Path("/swagger/json").Methods("GET").Handler(api.Handler())
    router.PathPrefix("/swagger/ui").Handler(swag.UIHandler("/swagger/ui", "/swagger/json", true))

    log.Fatal(http.ListenAndServe(":8080", router))
}
```

### echo

```go
func main() {
    ...
    
    router := echo.New()
    api.Walk(func (path string, e *swag.Endpoint) {
        h := echo.WrapHandler(e.Handler.(http.Handler))
        path = swag.ColonPath(path)
        
        switch strings.ToLower(e.Method) {
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
    router.GET("/swagger/ui/*", echo.WrapHandler(swag.UIHandler("/swagger/ui", "/swagger/json", true)))

    log.Fatal(http.ListenAndServe(":8080", router))
}
```

### httprouter

```go
func main() {
    ...
    
    router := httprouter.New()
    api.Walk(func (path string, e *swag.Endpoint) {
        h := e.Handler.(http.Handler)
        path = swag.ColonPath(path)
        router.Handler(e.Method, path, h)
    })
    
    router.Handler(http.MethodGet, "/swagger/json", api.Handler())
    router.Handler(http.MethodGet, "/swagger/ui/*any", swag.UIHandler("/swagger/ui", "/swagger/json", true))

    log.Fatal(http.ListenAndServe(":8080", router))
}
```

### fasthttp

```go
package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"path"
	"path/filepath"
	"strings"

	"github.com/fasthttp/router"
	"github.com/valyala/fasthttp"

	"github.com/zc2638/swag"
	"github.com/zc2638/swag/asserts"
	"github.com/zc2638/swag/endpoint"
	"github.com/zc2638/swag/option"
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

func handle(ctx *fasthttp.RequestCtx) {
	str := fmt.Sprintf("[%s] Hello World!", string(ctx.Method()))
	_, _ = ctx.Write([]byte(str))
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

	api := swag.New(
		option.Title("Example API Doc"),
		option.Security("petstore_auth", "read:pets"),
		option.SecurityScheme("petstore_auth",
			option.OAuth2Security("accessCode", "http://example.com/oauth/authorize", "http://example.com/oauth/token"),
			option.OAuth2Scope("write:pets", "modify pets in your account"),
			option.OAuth2Scope("read:pets", "read your pets"),
		),
		option.Endpoints(post, get),
	)
	api.AddEndpoint(test)

	r := router.New()
	api.Walk(func(path string, e *swag.Endpoint) {
		if v, ok := e.Handler.(func(ctx *fasthttp.RequestCtx)); ok {
			r.Handle(e.Method, path, fasthttp.RequestHandler(v))
		} else {
			r.Handle(e.Method, path, e.Handler.(fasthttp.RequestHandler))
		}
	})

	buildSchemeFn := func(ctx *fasthttp.RequestCtx) string {
		var scheme []byte

		if ctx.IsTLS() {
			scheme = []byte("https")
		}
		if v := ctx.Request.Header.Peek("X-Forwarded-Proto"); v != nil {
			scheme = v
		}
		if string(scheme) == "" {
			scheme = ctx.URI().Scheme()
		}
		if string(scheme) == "" {
			scheme = []byte("http")
		}
		return string(scheme)
	}

	doc := api.Clone()
	r.GET("/swagger/json", func(ctx *fasthttp.RequestCtx) {
		scheme := buildSchemeFn(ctx)
		doc.Host = string(ctx.Host())
		doc.Schemes = []string{scheme}

		b, err := json.Marshal(doc)
		if err != nil {
			ctx.Error("Parse API Doc exceptions", http.StatusInternalServerError)
			return
		}
		_, _ = ctx.Write(b)
	})

	r.ANY("/swagger/ui/{any:*}", func(ctx *fasthttp.RequestCtx) {
		currentPath := strings.TrimPrefix(string(ctx.Path()), "/swagger/ui")

		if currentPath == "/" || currentPath == "index.html" {
			fullName := path.Join(asserts.DistDir, "index.html")
			fileData, err := asserts.Dist.ReadFile(fullName)
			if err != nil {
				ctx.Error("index.html read exception", http.StatusInternalServerError)
				return
			}

			scheme := buildSchemeFn(ctx)
			currentURI := scheme + "://" + path.Join(string(ctx.Host()), "/swagger/json")

			fileData = bytes.ReplaceAll(fileData, []byte(asserts.URL), []byte(currentURI))
			ctx.SetContentType("text/html; charset=utf-8")
			ctx.Write(fileData)
			return
		}
		sfs := swag.DirFS(asserts.DistDir, asserts.Dist)
		file, err := sfs.Open(currentPath)
		if err != nil {
			ctx.Error(err.Error(), http.StatusInternalServerError)
			return
		}

		stat, err := file.Stat()
		if err != nil {
			ctx.Error(err.Error(), http.StatusInternalServerError)
			return
		}

		switch strings.TrimPrefix(filepath.Ext(stat.Name()), ".") {
		case "css":
			ctx.SetContentType("text/css; charset=utf-8")
		case "js":
			ctx.SetContentType("application/javascript")
		}
		io.Copy(ctx, file)
	})

	fasthttp.ListenAndServe(":8080", r.Handler)
}
```
