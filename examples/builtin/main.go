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
package main

import (
	"io"
	"net/http"

	"github.com/zc2638/swag"
	"github.com/zc2638/swag/endpoint"
	"github.com/zc2638/swag/swagger"
)

func handle(w http.ResponseWriter, req *http.Request) {
	io.WriteString(w, req.Method+" - Insert your code here")
}

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
	Test
}

type Test struct {
	A string `json:"a"`
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
		swag.Endpoints(post, get, test),
		swag.Security("petstore_auth", "read:pets"),
		swag.SecurityScheme("petstore_auth",
			swagger.OAuth2Security("accessCode", "http://example.com/oauth/authorize", "http://example.com/oauth/token"),
			swagger.OAuth2Scope("write:pets", "modify pets in your account"),
			swagger.OAuth2Scope("read:pets", "read your pets"),
		),
	)
	api.RegisterMuxWithData(http.DefaultServeMux, true)
	http.ListenAndServe(":8080", nil)
}
