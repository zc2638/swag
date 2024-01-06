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
	"encoding/json"
	"fmt"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

type Person struct {
	First string
}

type Anonymous struct {
	AnyOne string
}

type Pet struct {
	Friend          Person    `json:"friend" desc:"description short expression"`
	Friends         []Person  `json:"friends" description:"long desc"`
	Pointer         *Person   `json:"pointer" required:"true"`
	Pointers        []*Person `json:"pointers"`
	Int             int
	IntArray        []int
	Int64Array      []int64
	String          string
	StringSecondWay string `json:"StringSecondWay,string"`
	StringArray     []string
	Float           float32
	FloatArray      []float32
	Double          float64
	DoubleArray     []float64
	Bool            bool
	Enum            string `json:"enum" enum:"a,b,c" example:"b"`
	Anonymous
	MapSlicePtr       map[string][]*string
	MapSlice          map[string][]string
	MapSliceStructPtr map[string][]*Person
	MapSliceStruct    map[string][]Person
	SliceStructPtr    *[]*Person
	SliceStruct       *[]Person
	SliceStringPtr    *[]*string
	SliceString       *[]string
	MapNestOptions    *MapObj `json:"map_nest_options,omitempty"`
}

type MapObj struct {
	RuleOptions map[string]*MapOption `json:"rule_options"`
}

type MapOption struct {
	Name       string                `json:"name"`
	SubOptions map[string]*MapOption `json:"sub_options,omitempty"`
}

type Empty struct {
	Nope int `json:"-"`
}

func TestDefine(t *testing.T) {
	v := define(Pet{})
	obj, ok := v["swag.Pet"]
	assert.True(t, ok)
	assert.False(t, obj.IsArray)
	assert.Equal(t, 26, len(obj.Properties))

	content := make(map[string]Object)
	data, err := os.ReadFile("testdata/pet.json")
	assert.Nil(t, err)
	err = json.NewDecoder(bytes.NewReader(data)).Decode(&content)
	assert.Nil(t, err)
	expected := content["swag.Pet"]

	assert.Equal(t, expected.IsArray, obj.IsArray, "expected IsArray to match")
	assert.Equal(t, expected.Type, obj.Type, "expected Type to match")
	assert.Equal(t, expected.Required, obj.Required, "expected Required to match")
	assert.Equal(t, len(expected.Properties), len(obj.Properties), "expected same number of properties")

	for k, v := range obj.Properties {
		e := expected.Properties[k]
		assert.Equal(t, e.Type, v.Type, "expected %v.Type to match", k)
		assert.Equal(t, e.Description, v.Description, "expected %v.Required to match", k)
		assert.Equal(t, e.Enum, v.Enum, "expected %v.Required to match", k)
		assert.Equal(t, e.Format, v.Format, "expected %v.Required to match", k)
		assert.Equal(t, e.Ref, v.Ref, "expected %v.Required to match", k)
		assert.Equal(t, e.Example, v.Example, "expected %v.Required to match", k)
		assert.Equal(t, e.Items, v.Items, "expected %v.Required to match", k)
	}
}

func TestNotStructDefine(t *testing.T) {
	v := define(int32(1))
	obj, ok := v["int32"]
	assert.True(t, ok)
	assert.False(t, obj.IsArray)
	assert.Equal(t, "integer", obj.Type)
	assert.Equal(t, "int32", obj.Format)

	v = define(uint64(1))
	obj, ok = v["uint64"]
	assert.True(t, ok)
	assert.False(t, obj.IsArray)
	assert.Equal(t, "integer", obj.Type)
	assert.Equal(t, "int64", obj.Format)

	v = define("")
	obj, ok = v["string"]
	assert.True(t, ok)
	assert.False(t, obj.IsArray)
	assert.Equal(t, "string", obj.Type)
	assert.Equal(t, "", obj.Format)

	v = define(byte(1))
	obj, ok = v["uint8"]
	if !assert.True(t, ok) {
		fmt.Printf("%v", v)
	}
	assert.False(t, obj.IsArray)
	assert.Equal(t, "integer", obj.Type)
	assert.Equal(t, "int32", obj.Format)

	v = define([]byte{1, 2})
	obj, ok = v["uint8"]
	if !assert.True(t, ok) {
		fmt.Printf("%v", v)
	}
	assert.True(t, obj.IsArray)
	assert.Equal(t, "integer", obj.Type)
	assert.Equal(t, "int32", obj.Format)
}

func TestHonorJsonIgnore(t *testing.T) {
	v := define(Empty{})
	obj, ok := v["swag.Empty"]
	assert.True(t, ok)
	assert.False(t, obj.IsArray)
	assert.Equal(t, 0, len(obj.Properties), "expected zero exposed properties")
}

func TestMakeSchemaType(t *testing.T) {
	sliceSchema := MakeSchema([]string{})
	assert.Equal(t, "array", sliceSchema.Type, "expect array type but get %s", sliceSchema.Type)

	objSchema := MakeSchema(struct{}{})
	assert.Equal(t, "", objSchema.Type, "expect array type but get %s", objSchema.Type)
}
