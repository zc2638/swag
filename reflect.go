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
	"github.com/zc2638/swag/types"
	"reflect"
	"strings"
)

func inspect(t reflect.Type, jsonTag string) Property {
	p := Property{
		GoType: t,
	}

	if strings.Contains(jsonTag, ",string") {
		p.Type = "string"
		return p
	}

	if p.GoType.Kind() == reflect.Ptr {
		p.GoType = p.GoType.Elem()
	}

	switch p.GoType.Kind() {
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Uint8, reflect.Uint16, reflect.Uint32:
		p.Type = types.Integer.String()
		p.Format = "int32"

	case reflect.Int64, reflect.Uint64:
		p.Type = types.Integer.String()
		p.Format = "int64"

	case reflect.Float64:
		p.Type = types.Number.String()
		p.Format = "double"

	case reflect.Float32:
		p.Type = types.Number.String()
		p.Format = "float"

	case reflect.Bool:
		p.Type = types.Boolean.String()

	case reflect.String:
		p.Type = types.String.String()

	case reflect.Struct:
		name := makeName(p.GoType)
		p.Ref = makeRef(name)

	case reflect.Ptr:
		p.GoType = t.Elem()
		name := makeName(p.GoType)
		p.Ref = makeRef(name)

	case reflect.Map:
		p.Type = "object"
		p.AddPropertie = buildMapType(p.GoType)
		// get the go reflect type of the map value
		p.GoType = p.AddPropertie.GoType

	case reflect.Slice:
		p.Type = types.Array.String()
		p.Items = &Items{}

		// dereference the slice,get the element object of the slice
		// Example:
		//    t ==> *[]*Struct|Struct|*string|string
		//    p.GoType.Kind() ==> []*Struct|Struct|*string|string
		if p.GoType.Kind() == reflect.Slice {
			//p.GoType ==> *Struct|Struct|*string|string
			p.GoType = p.GoType.Elem()
		}

		switch p.GoType.Kind() {
		case reflect.Ptr:
			//p.GoType ==> Struct|string
			p.GoType = p.GoType.Elem()

			// determine the type of element object in the slice
			isPrimitive := isPrimitiveType(p.GoType.Name(), p.GoType)

			if isPrimitive {
				// golang built-in primitive type
				kind_type := jsonSchemaType(p.GoType.String(), p.GoType)
				p.Items.Type = kind_type
			} else {
				// Struct types
				name := makeName(p.GoType)
				p.Items.Ref = makeRef(name)
			}

		case reflect.Struct:
			name := makeName(p.GoType)
			p.Items.Ref = makeRef(name)

		case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Uint8, reflect.Uint16, reflect.Uint32:
			p.Items.Type = types.Integer.String()
			p.Items.Format = "int32"

		case reflect.Int64, reflect.Uint64:
			p.Items.Type = types.Integer.String()
			p.Items.Format = "int64"

		case reflect.Float64:
			p.Items.Type = types.Number.String()
			p.Items.Format = "double"

		case reflect.Float32:
			p.Items.Type = types.Number.String()
			p.Items.Format = "float"

		case reflect.String:
			p.Items.Type = types.String.String()
		}
	}

	return p
}

// buildMapType build map type swag info
func buildMapType(mapType reflect.Type) *AdditionalProperties {

	prop := AdditionalProperties{}

	// get the element object of the map
	// Example:
	//	mapType ==> map[string][]*Struct|Struct|*string|string
	//  or mapTYpe ==> map[string]*Struct|Struct|*string|string
	//	mapType.Elem().Kind() ==> []*Struct|Struct|*string|string
	//  or mapType.Elem().Kind() ==> *Struct|Struct|*string|string
	if mapType.Elem().Kind().String() != "interface" {
		isSlice := isSliceOrArryType(mapType.Elem().Kind())
		if isSlice || isByteArrayType(mapType.Elem()) {
			// if map value is slice
			// Example:
			//   mapType.Elem()==> []*Struct|Struct|*string|string
			mapType = mapType.Elem()
		}

		// if map value is struct or built-in primitive type
		// Example:
		//    mapType.Elem()==> *Struct|Struct|*string|string
		isPrimitive := isPrimitiveType(mapType.Elem().Name(), mapType.Elem())

		if isByteArrayType(mapType.Elem()) {
			prop.Type = "string"
		} else {
			if isSlice {
				prop.Type = types.Array.String()
				prop.Items = &Items{}
				if isPrimitive {
					prop.Items.Type = jsonSchemaType(mapType.Elem().String(), mapType.Elem())
					prop.GoType = getGoType(mapType.Elem())
				} else {
					prop.Items.Ref = makeMapRef(mapType.Elem().String())
					prop.GoType = getGoType(mapType.Elem())
				}
			} else if isPrimitive {
				prop.Type = jsonSchemaType(mapType.Elem().String(), mapType.Elem())
				prop.GoType = getGoType(mapType.Elem())
			} else {
				prop.Ref = makeMapRef(mapType.Elem().String())
				prop.GoType = getGoType(mapType.Elem())
			}
		}
	}

	return &prop
}

func getGoType(t reflect.Type) reflect.Type {

	var goType reflect.Type

	if t.Kind() == reflect.Ptr {
		goType = t.Elem()
	} else {
		goType = t
	}

	return goType
}

func makeMapRef(typeName string) string {
	type_name := strings.Trim(typeName, "*")
	return makeRef(type_name)
}

func isSliceOrArryType(t reflect.Kind) bool {
	return t == reflect.Slice || t == reflect.Array
}

func isByteArrayType(t reflect.Type) bool {
	return (t.Kind() == reflect.Slice || t.Kind() == reflect.Array) &&
		t.Elem().Kind() == reflect.Uint8
}

// isPrimitiveType Whether it is a built-in primitive type
func isPrimitiveType(modelName string, modelType reflect.Type) bool {
	var modelKind reflect.Kind

	if modelType.Kind() == reflect.Ptr {
		modelKind = modelType.Elem().Kind()
	} else {
		modelKind = modelType.Kind()
	}

	switch modelKind {
	case reflect.Bool:
		return true
	case reflect.Float32, reflect.Float64,
		reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64,
		reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		return true
	case reflect.String:
		return true
	}

	if len(modelName) == 0 {
		return false
	}

	return strings.Contains("time.Time time.Duration json.Number", modelName)
}

func jsonSchemaType(modelName string, modelType reflect.Type) string {
	var modelKind reflect.Kind

	if modelType.Kind() == reflect.Ptr {
		modelKind = modelType.Elem().Kind()
	} else {
		modelKind = modelType.Kind()
	}

	schemaMap := map[string]string{
		"time.Time":     "string",
		"time.Duration": "integer",
		"json.Number":   "number",
	}

	if mapped, ok := schemaMap[modelName]; ok {
		return mapped
	}

	// check if original type is primitive
	switch modelKind {
	case reflect.Bool:
		return "boolean"
	case reflect.Float32, reflect.Float64:
		return "number"
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64,
		reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		return "integer"
	case reflect.String:
		return "string"
	}

	return modelName // use as is (custom or struct)
}

func buildProperty(t reflect.Type) (map[string]Property, []string) {
	properties := make(map[string]Property)
	required := make([]string, 0)
	for i := 0; i < t.NumField(); i++ {
		field := t.Field(i)

		// skip unexported fields
		if strings.ToLower(field.Name[0:1]) == field.Name[0:1] {
			continue
		}
		if field.Anonymous {
			// 暂不处理匿名结构的required
			ps, _ := buildProperty(field.Type)
			for name, p := range ps {
				properties[name] = p
			}
			continue
		}

		// determine the json name of the field
		name := strings.TrimSpace(field.Tag.Get("json"))
		if name == "" || strings.HasPrefix(name, ",") {
			name = field.Name

		} else {
			// strip out things like , omitempty
			parts := strings.Split(name, ",")
			name = parts[0]
		}

		parts := strings.Split(name, ",") // foo,omitempty => foo
		name = parts[0]
		if name == "-" {
			// honor json ignore tag
			continue
		}
		p := inspect(field.Type, field.Tag.Get("json"))

		// determine the extra info of the field
		if _, ok := field.Tag.Lookup("required"); ok {
			required = append(required, name)
		}
		if example := field.Tag.Get("example"); example != "" {
			p.Example = example
		}
		if description := field.Tag.Get("description"); description != "" {
			p.Description = description
		}
		if desc := field.Tag.Get("desc"); desc != "" {
			p.Description = desc
		}
		if enum := field.Tag.Get("enum"); enum != "" {
			p.Enum = strings.Split(enum, ",")
		}
		properties[name] = p
	}
	return properties, required
}

func defineObject(v interface{}, desc string) Object {
	var t reflect.Type
	switch value := v.(type) {
	case reflect.Type:
		t = value
	default:
		t = reflect.TypeOf(v)
	}

	isArray := t.Kind() == reflect.Slice
	if isArray {
		t = t.Elem()
	}
	if t.Kind() == reflect.Ptr {
		t = t.Elem()
	}

	if t.Kind() != reflect.Struct {
		p := inspect(t, "")
		return Object{
			IsArray: isArray,
			GoType:  t,
			Type:    p.Type,
			Format:  p.Format,
			Name:    t.Kind().String(),
		}
	}
	properties, required := buildProperty(t)

	return Object{
		IsArray:     isArray,
		GoType:      t,
		Type:        "object",
		Name:        makeName(t),
		Required:    required,
		Properties:  properties,
		Description: desc,
	}
}

func define(v interface{}) map[string]Object {
	objMap := map[string]Object{}

	obj := defineObject(v, "")
	objMap[obj.Name] = obj

	dirty := true

	for dirty {
		dirty = false
		for _, d := range objMap {
			for _, p := range d.Properties {
				if p.GoType.Kind() == reflect.Struct {
					name := makeName(p.GoType)
					if _, exists := objMap[name]; !exists {
						child := defineObject(p.GoType, p.Description)
						objMap[child.Name] = child
						dirty = true
					}
				}
			}
		}
	}

	return objMap
}

// MakeSchema takes struct or pointer to a struct and returns a Schema instance suitable for use by the swagger doc
func MakeSchema(prototype interface{}) *Schema {
	schema := &Schema{
		Prototype: prototype,
	}

	obj := defineObject(prototype, "")
	if obj.IsArray {
		schema.Type = "array"
		schema.Items = &Items{
			Ref: makeRef(obj.Name),
		}

	} else {
		schema.Ref = makeRef(obj.Name)
	}

	return schema
}
