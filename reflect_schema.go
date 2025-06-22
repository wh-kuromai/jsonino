package jsonino

import (
	"errors"
	"reflect"
)

// StructToSchema converts a Go struct to an ObjectScheme, handling nested structs recursively.
func StructToSchema(v interface{}) (SchemaNode, error) {

	var typ reflect.Type
	typ, ok := v.(reflect.Type)
	if !ok {
		typ = reflect.TypeOf(v)
	}

	if typ.Kind() == reflect.Ptr {
		typ = typ.Elem()
	}
	if typ.Kind() != reflect.Struct {
		return nil, errors.New("StructToSchema: expected struct type")
	}

	//fmt.Println("--", typ.Name())

	schema := &ObjectScheme{
		TypeName:   "object",
		Properties: map[string]*SchemaNodeBase{},
		Required:   []string{},
	}

	for i := 0; i < typ.NumField(); i++ {
		field := typ.Field(i)
		jsonTag := field.Tag.Get("json")
		if jsonTag == "-" {
			continue
		}

		name := field.Name
		if jsonTag != "" {
			commaIdx := indexComma(jsonTag)
			if commaIdx >= 0 {
				name = jsonTag[:commaIdx]
			} else {
				name = jsonTag
			}
		}

		schemaNode, err := TypeToSchema(reflect.New(field.Type).Elem().Type())
		if err != nil {
			return nil, err
		}

		base := &SchemaNodeBase{
			TypeName: schemaNode.Type(),
			This:     schemaNode,
		}

		if required := field.Tag.Get("required"); required == "true" {
			schema.Required = append(schema.Required, name)
		}

		schema.Properties[name] = base
	}

	return schema, nil
}

// TypeToSchema converts a Go value into a SchemaNode.
func TypeToSchema(v any) (SchemaNode, error) {
	if v == nil {
		return nil, nil
	}

	var typ reflect.Type
	typ, ok := v.(reflect.Type)
	if !ok {
		typ = reflect.TypeOf(v)
	}

	if typ.Kind() == reflect.Ptr {
		typ = typ.Elem()
	}

	// Special case: time.Time -> string
	if typ.PkgPath() == "time" && typ.Name() == "Time" {
		return &StringScheme{TypeName: "string"}, nil
	}

	switch typ.Kind() {
	case reflect.String:
		return &StringScheme{TypeName: "string"}, nil
	case reflect.Bool:
		return &StringScheme{TypeName: "boolean"}, nil // adjust if BooleanScheme is defined
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64,
		reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64,
		reflect.Float32, reflect.Float64:
		return &NumberScheme{TypeName: "number"}, nil
	case reflect.Slice, reflect.Array:
		itemType := typ.Elem()
		itemNode, err := TypeToSchema(reflect.New(itemType).Elem().Interface())
		if err != nil {
			return nil, err
		}
		return &ArrayScheme{
			TypeName: "array",
			Items: &SchemaNodeBase{
				TypeName: itemNode.Type(),
				This:     itemNode,
			},
		}, nil
	case reflect.Struct:
		return StructToSchema(v)
	default:
		return nil, errors.New("unsupported type: " + typ.String())
	}
}

func indexComma(s string) int {
	for i := 0; i < len(s); i++ {
		if s[i] == ',' {
			return i
		}
	}
	return -1
}

func goKindToTypeName(k reflect.Kind) string {
	switch k {
	case reflect.String:
		return "string"
	case reflect.Int, reflect.Float64, reflect.Int64:
		return "number"
	case reflect.Bool:
		return "boolean"
	default:
		return "unknown"
	}
}
