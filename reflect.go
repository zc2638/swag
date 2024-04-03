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
	"reflect"
	"strings"

	"github.com/zc2638/swag/types"
)

func inspect(t reflect.Type) Property {
	// unwrap pointer
	if t.Kind() == reflect.Ptr {
		return inspect(t.Elem())
	}

	p := Property{GoType: t}

	p.Type = types.Parse(p.GoType).String()
	p.Format = types.ParseFormat(p.GoType).String()
	switch p.GoType.Kind() {
	case reflect.Struct:
		name := makeName(p.GoType)
		p.Ref = makeRef(name)
		p.Type = ""
	case reflect.Map:
		ap := inspect(p.GoType.Elem())
		p.AdditionalProperties = &ap
	case reflect.Slice, reflect.Array:
		// unwrap array or slice
		p.GoType = p.GoType.Elem()
		items := inspect(p.GoType)
		p.Items = &items
	default:
	}
	return p
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

		var p Property
		jsonTag := field.Tag.Get("json")
		if strings.Contains(jsonTag, ",string") {
			p.Type = types.String.String()
			p.GoType = field.Type
		} else {
			p = inspect(field.Type)
		}

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

func defineObject(v interface{}) Object {
	var t reflect.Type
	switch value := v.(type) {
	case reflect.Type:
		t = value
	default:
		t = reflect.TypeOf(v)
	}

	kind := t.Kind()
	if kind == reflect.Ptr {
		// unwrap pointer
		return defineObject(t.Elem())
	}

	obj := Object{
		Name: kind.String(),
	}
	switch kind {
	case reflect.Slice, reflect.Array:
		obj.Name = makeName(t)
		t = t.Elem()
		p := inspect(t)
		obj.Items = &p
		obj.Type = types.Array.String()
	case reflect.Map:
		obj.Name = makeName(t)
		t = t.Elem()
		p := inspect(t)
		obj.AdditionalProperties = &p
		obj.Type = types.Object.String()
	case reflect.Struct:
		obj.Name = makeName(t)
		properties, required := buildProperty(t)
		obj.Type = types.Object.String()
		obj.Required = required
		obj.Properties = properties
	default:
		p := inspect(t)
		obj.Type = p.Type
		obj.Format = p.Format
	}
	return obj
}

func define(v interface{}) map[string]Object {
	objMap := make(map[string]Object)
	obj := defineObject(v)
	objMap[obj.Name] = obj

	dirty := true
	for dirty {
		dirty = false
		for _, d := range objMap {
			for _, p := range d.Properties {
				prop := &p
				for prop != nil {
					if prop.GoType.Kind() == reflect.Struct {
						name := makeName(prop.GoType)
						if _, exists := objMap[name]; !exists {
							child := defineObject(prop.GoType)
							objMap[child.Name] = child
							dirty = true
						}
					}
					prop = prop.AdditionalProperties
				}
			}
		}
	}
	return objMap
}

// MakeSchema takes struct or pointer to a struct and returns a Schema instance suitable for use by the swagger doc
func MakeSchema(prototype interface{}) *Schema {
	obj := defineObject(prototype)
	schema := &Schema{
		Prototype: prototype,
		Ref:       makeRef(obj.Name),
	}
	return schema
}
