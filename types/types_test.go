// Package types
// Created by zc on 2022/6/11.
package types

import (
	"context"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestParameterType(t *testing.T) {
	assert.Equal(t, "integer", Integer.String())
	assert.Equal(t, "number", Number.String())
	assert.Equal(t, "boolean", Boolean.String())
	assert.Equal(t, "string", String.String())
	assert.Equal(t, "array", Array.String())
	assert.Equal(t, "file", File.String())
}

func TestParameterType_String(t *testing.T) {
	tests := []struct {
		name string
		pt   ParameterType
		want string
	}{
		{
			name: "integer",
			pt:   Integer,
			want: "integer",
		},
		{
			name: "number",
			pt:   Number,
			want: "number",
		},
		{
			name: "boolean",
			pt:   Boolean,
			want: "boolean",
		},
		{
			name: "string",
			pt:   String,
			want: "string",
		},
		{
			name: "array",
			pt:   Array,
			want: "array",
		},
		{
			name: "file",
			pt:   File,
			want: "file",
		},
		{
			name: "unknown",
			pt:   ParameterType("unknown"),
			want: "unknown",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equalf(t, tt.want, tt.pt.String(), "String()")
		})
	}
}

func TestAddURLParamsToContext(t *testing.T) {
	params1 := map[string]string{
		"a": "1",
		"b": "2",
	}
	params2 := map[string]string{
		"c": "3",
		"d": "4",
	}

	type args struct {
		parent context.Context
		params map[string]string
	}
	tests := []struct {
		name string
		args args
		want map[string]string
	}{
		{
			name: "default",
			args: args{
				parent: context.Background(),
				params: params1,
			},
			want: params1,
		},
		{
			name: "add",
			args: args{
				parent: AddURLParamsToContext(context.Background(), params1),
				params: params2,
			},
			want: map[string]string{
				"a": "1",
				"b": "2",
				"c": "3",
				"d": "4",
			},
		},
		{
			name: "cover",
			args: args{
				parent: AddURLParamsToContext(context.Background(), params1),
				params: map[string]string{
					"a": "10",
				},
			},
			want: map[string]string{
				"a": "10",
				"b": "2",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			currentCtx := AddURLParamsToContext(tt.args.parent, tt.args.params)
			routeVal := currentCtx.Value(RouteContextKey)
			routeCtx, ok := routeVal.(*Context)
			if !ok {
				t.Error("AddURLParamsToContext failed: can not get RouteContext")
				return
			}
			assert.Equalf(t, tt.want, routeCtx.PathParams, "AddURLParamsToContext(%v, %v)", tt.args.parent, tt.args.params)
		})
	}
}

func TestURLParam(t *testing.T) {
	req, _ := http.NewRequest(http.MethodGet, "/", nil)
	ctx := AddURLParamsToContext(req.Context(), map[string]string{
		"name": "zc",
	})
	req = req.WithContext(ctx)

	type args struct {
		r   *http.Request
		key string
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			name: "no ctx",
			args: args{
				r:   new(http.Request),
				key: "none",
			},
			want: "",
		},
		{
			name: "no value",
			args: args{
				r:   req,
				key: "none",
			},
			want: "",
		},
		{
			name: "value exists",
			args: args{
				r:   req,
				key: "name",
			},
			want: "zc",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equalf(t, tt.want, URLParam(tt.args.r, tt.args.key), "URLParam(%v, %v)", tt.args.r, tt.args.key)
		})
	}
}

func TestURLParamFromCtx(t *testing.T) {
	ctx := AddURLParamsToContext(context.Background(), map[string]string{
		"name": "zc",
	})

	type args struct {
		ctx context.Context
		key string
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			name: "no ctx",
			args: args{
				ctx: nil,
				key: "none",
			},
			want: "",
		},
		{
			name: "no value",
			args: args{
				ctx: ctx,
				key: "none",
			},
			want: "",
		},
		{
			name: "value exists",
			args: args{
				ctx: ctx,
				key: "name",
			},
			want: "zc",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equalf(t, tt.want, URLParamFromCtx(tt.args.ctx, tt.args.key), "URLParamFromCtx(%v, %v)", tt.args.ctx, tt.args.key)
		})
	}
}
