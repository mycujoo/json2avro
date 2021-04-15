package main

import (
	"encoding/json"
	"log"
	"reflect"
	"strconv"
	"strings"
)

func Parse(title string, payload []byte, nullables bool) ([]byte, error) {
	var obj map[string]interface{}

	if err := json.Unmarshal(payload, &obj); err != nil {
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
		case float64:
			var typ string

			stringified := strconv.FormatFloat(val, 'E', -1, 64)
			if strings.Contains(stringified, ".") {
				if float64(float32(val)) == val {
					typ = "float"
				} else {
					typ = "double"
				}
			} else {
				if float64(int32(val)) == val {
					typ = "int"
				} else if float64(int64(val)) == val {
					typ = "long"
				}
			}

			schema = append(schema, map[string]interface{}{
				"name": k,
				"type": _type(typ, nullables),
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
