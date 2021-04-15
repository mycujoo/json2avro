package main

import (
	"encoding/json"
	"io"
	"log"
	"reflect"
	"strings"
)

func Parse(title string, payload io.Reader, nullables bool) ([]byte, error) {
	var obj map[string]interface{}

	decoder := json.NewDecoder(payload)
	decoder.UseNumber()

	if err := decoder.Decode(&obj); err != nil {
		return nil, err
	}

	uniques := make(map[string]interface{})

	fields := parseObject(uniques, obj, nullables)

	return json.Marshal(map[string]interface{}{
		"name":   strings.ReplaceAll(title, " ", ""),
		"type":   "record",
		"fields": fields,
	})
}

func parseObject(uniques map[string]interface{}, obj map[string]interface{}, nullables bool) []map[string]interface{} {
	schema := make([]map[string]interface{}, 0, len(obj))

	for k, v := range obj {
		switch val := v.(type) {
		case string:
			schema = append(schema, map[string]interface{}{
				"name": k,
				"type": _type("string", nullables),
			})
		case bool:
			schema = append(schema, map[string]interface{}{
				"name": k,
				"type": _type("boolean", nullables),
			})
		case json.Number:
			if _, err := val.Int64(); err == nil {
				schema = append(schema, map[string]interface{}{
					"name": k,
					"type": _type("long", nullables),
				})
				continue
			}

			schema = append(schema, map[string]interface{}{
				"name": k,
				"type": _type("double", nullables),
			})
		case []interface{}:
			var typ interface{}

			for _, y := range val {
				switch y.(type) {
				case int:
					typ = "int"
				case string:
					typ = "string"
				case map[string]interface{}:
					parsed, ok := y.(map[string]interface{})
					if !ok {
						log.Fatalf("%s.%s is not map[string]interface{}", k, y)
					}
					typ = map[string]interface{}{
						"name":   _name(uniques, k+"_record"),
						"type":   "record",
						"fields": parseObject(uniques, parsed, nullables),
					}
				}
			}

			schema = append(schema, map[string]interface{}{
				"name":    k,
				"default": []interface{}{},
				"type": _type(map[string]interface{}{
					"name":  _name(uniques, k+"_array"),
					"type":  "array",
					"items": typ,
				}, nullables),
			})
		case map[string]interface{}:
			schema = append(schema, map[string]interface{}{
				"name": k,
				"type": _type(map[string]interface{}{
					"name":   _name(uniques, k+"_record"),
					"type":   "record",
					"fields": parseObject(uniques, val, nullables),
				}, nullables),
			})
		default:
			log.Println(100, k, reflect.TypeOf(v))
		}
	}

	return schema
}

func _name(uniques map[string]interface{}, name string) string {
	name = strings.ToLower(name)
	for {
		_, ok := uniques[name]
		if !ok {
			break
		}
		name = name + "_"
	}
	uniques[name] = 1
	return name
}

func _type(typ interface{}, nullables bool) interface{} {
	if nullables {
		return []interface{}{typ, "null"}
	}

	return typ
}
