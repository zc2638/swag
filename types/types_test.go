// Package types
// Created by zc on 2022/6/11.
package types

import (
	"context"
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
