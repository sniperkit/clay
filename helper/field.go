package helper

import (
	"errors"
	"reflect"
	"strings"

	"github.com/qb0C80aE/clay/extension"

	"github.com/serenize/snaker"
)

// AssociationType is the association type of GORM
type AssociationType int

const (
	none AssociationType = iota
	belongsTo
	hasMany
	hasOne
)

func contains(ss map[string]interface{}, s string) bool {
	_, ok := ss[s]

	return ok
}

func merge(m1, m2 map[string]interface{}) map[string]interface{} {
	result := make(map[string]interface{})

	for k, v := range m1 {
		result[k] = v
	}

	for k, v := range m2 {
		result[k] = v
	}

	return result
}

// QueryFields builds a query field string based on arguments
// It requires a container which has actual fields, not a model which might not have actual fields
func QueryFields(container interface{}, fields map[string]interface{}, targetTag string) (string, error) {
	if _, ok := container.(extension.Model); ok {
		return "", errors.New("given data is not a model, not a container")
	}

	var targetFieldKey string

	vs := reflect.ValueOf(container)

	for vs.Kind() == reflect.Ptr {
		vs = vs.Elem()
	}
	if !vs.IsValid() {
		return "", errors.New("given data is not valid")
	}

	assocs := make(map[string]AssociationType)

	ts := vs.Type()
	for i := 0; i < ts.NumField(); i++ {
		f := ts.Field(i)

		targetFieldTag := f.Tag.Get(targetTag)

		if targetFieldTag == "" {
			targetFieldKey = f.Name
		} else {
			targetFieldKey = strings.Split(targetFieldTag, ",")[0]
		}

		switch vs.Field(i).Kind() {
		case reflect.Ptr:
			if _, ok := ts.FieldByName(f.Name + "ID"); ok {
				assocs[targetFieldKey] = belongsTo
			} else {
				assocs[targetFieldKey] = hasOne
			}
		case reflect.Slice:
			assocs[targetFieldKey] = hasMany
		default:
			assocs[targetFieldKey] = none
		}
	}

	result := []string{}

	for k := range fields {
		if k == "*" {
			return "*", nil
		}

		if _, ok := assocs[k]; !ok {
			continue
		}

		switch assocs[k] {
		case none:
			result = append(result, k)
		case belongsTo:
			result = append(result, k+"_id")
		default:
			result = append(result, "id")
		}
	}

	return strings.Join(result, ","), nil
}

// ParseFields parses the argument and generate a map based on that
func ParseFields(fields string) map[string]interface{} {
	result := make(map[string]interface{})

	if fields == "*" {
		result["*"] = nil
		return result
	}

	for _, field := range strings.Split(fields, ",") {
		parts := strings.SplitN(field, ".", 2)

		if len(parts) == 2 {
			if result[parts[0]] == nil {
				result[parts[0]] = ParseFields(parts[1])
			} else {
				result[parts[0]] = merge(result[parts[0]].(map[string]interface{}), ParseFields(parts[1]))
			}
		} else {
			result[parts[0]] = nil
		}
	}

	return result
}

func isEmptyValue(v reflect.Value) bool {
	switch v.Kind() {
	case reflect.Array, reflect.Map, reflect.Slice, reflect.String:
		return v.Len() == 0
	case reflect.Bool:
		return !v.Bool()
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		return v.Int() == 0
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Uintptr:
		return v.Uint() == 0
	case reflect.Float32, reflect.Float64:
		return v.Float() == 0
	case reflect.Interface, reflect.Ptr:
		return v.IsNil()
	}
	return false
}

// FieldToMap generates a map from the model in argument
func FieldToMap(model interface{}, fields map[string]interface{}, targetTag string) (map[string]interface{}, error) {
	switch targetTag {
	case extension.TagYAML:
		// use as it is
	default:
		targetTag = extension.TagJSON
	}

	u := make(map[string]interface{})
	vs := reflect.ValueOf(model)

	for vs.Kind() == reflect.Ptr {
		vs = vs.Elem()
	}

	if !vs.IsValid() {
		return nil, errors.New("Invalid Parameter. The model is pointing nil")
	}

	if vs.Kind() != reflect.Struct {
		return nil, errors.New("Invalid Parameter. The specified parameter does not have a structure")
	}

	if !contains(fields, "*") {
		for field := range fields {
			if !vs.FieldByName(snaker.SnakeToCamel(field)).IsValid() {
				return nil, errors.New("Invalid Parameter. The specified field does not exist")
			}
		}
	}

	var fieldTagKey string
	var omitEmpty bool

	ts := vs.Type()
	for i := 0; i < ts.NumField(); i++ {
		field := ts.Field(i)
		jsonTag := field.Tag.Get(targetTag)
		omitEmpty = false

		if jsonTag == "" {
			fieldTagKey = field.Name
		} else {
			ss := strings.Split(jsonTag, ",")
			fieldTagKey = ss[0]

			if len(ss) > 1 && ss[1] == "omitempty" {
				omitEmpty = true
			}
		}

		if contains(fields, "*") {
			if !omitEmpty || !isEmptyValue(vs.Field(i)) {
				u[fieldTagKey] = vs.Field(i).Interface()
			}

			continue
		}

		if contains(fields, fieldTagKey) {
			v := fields[fieldTagKey]

			if vs.Field(i).Kind() == reflect.Ptr {
				if !vs.Field(i).IsNil() {
					if v == nil {
						u[fieldTagKey] = vs.Field(i).Elem().Interface()
					} else {
						k, err := FieldToMap(vs.Field(i).Elem().Interface(), v.(map[string]interface{}), targetTag)

						if err != nil {
							return nil, err
						}

						u[fieldTagKey] = k
					}
				} else {
					if v == nil {
						u[fieldTagKey] = nil
					} else {
						return nil, errors.New("Invalid Parameter. The structure is null")
					}
				}
			} else if vs.Field(i).Kind() == reflect.Slice {
				var fieldMap []interface{}
				s := reflect.ValueOf(vs.Field(i).Interface())

				for i := 0; i < s.Len(); i++ {
					if v == nil {
						fieldMap = append(fieldMap, s.Index(i).Interface())
					} else {

						if s.Index(i).Kind() == reflect.Ptr {
							k, err := FieldToMap(s.Index(i).Elem().Interface(), v.(map[string]interface{}), targetTag)

							if err != nil {
								return nil, err
							}

							fieldMap = append(fieldMap, k)
						} else {
							k, err := FieldToMap(s.Index(i).Interface(), v.(map[string]interface{}), targetTag)

							if err != nil {
								return nil, err
							}

							fieldMap = append(fieldMap, k)
						}
					}
				}

				u[fieldTagKey] = fieldMap
			} else {
				if v == nil {
					u[fieldTagKey] = vs.Field(i).Interface()
				} else {
					k, err := FieldToMap(vs.Field(i).Interface(), v.(map[string]interface{}), targetTag)

					if err != nil {
						return nil, err
					}

					u[fieldTagKey] = k
				}
			}
		}
	}

	return u, nil
}
