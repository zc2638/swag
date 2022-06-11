// Package types
// Created by zc on 2022/6/11.
package types

import (
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
