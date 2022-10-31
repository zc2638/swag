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

package types

import (
	"context"
	"net/http"
)

type ParameterType string

func (pt ParameterType) String() string {
	return string(pt)
}

const (
	Integer ParameterType = "integer"
	Number  ParameterType = "number"
	Boolean ParameterType = "boolean"
	String  ParameterType = "string"
	Array   ParameterType = "array"

	File ParameterType = "file"
)

type contextKey int

const (
	RouteContextKey contextKey = iota
)

type Context struct {
	PathParams map[string]string
}

// AddURLParamsToContext returns a copy of parent in which the context value is set
func AddURLParamsToContext(parent context.Context, params map[string]string) context.Context {
	routeVal := parent.Value(RouteContextKey)
	routeCtx, ok := routeVal.(*Context)
	if !ok {
		routeCtx = &Context{}
	}
	if routeCtx.PathParams == nil {
		routeCtx.PathParams = make(map[string]string)
	}
	for k, v := range params {
		routeCtx.PathParams[k] = v
	}
	return context.WithValue(parent, RouteContextKey, routeCtx)
}

// URLParam returns the url parameter from a http.Request object.
func URLParam(r *http.Request, key string) string {
	return URLParamFromCtx(r.Context(), key)
}

// URLParamFromCtx returns the url parameter from a http.Request Context.
func URLParamFromCtx(ctx context.Context, key string) string {
	if ctx == nil {
		return ""
	}

	routeVal := ctx.Value(RouteContextKey)
	routeCtx, ok := routeVal.(*Context)
	if !ok {
		return ""
	}
	for k, v := range routeCtx.PathParams {
		if k == key {
			return v
		}
	}
	return ""
}
