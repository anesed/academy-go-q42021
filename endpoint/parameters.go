package endpoint

import (
	"encoding/json"
	"errors"
	"net/http"
	"reflect"
	"strconv"
	"strings"
)

func WrapHandler(handler http.Handler) http.HandlerFunc {
	fields := reflect.ValueOf(handler.(interface{})).Elem()
	reflectType := reflect.Indirect(fields).Type()
	preprocess := func(writer http.ResponseWriter, r *http.Request, errorBag []error) []error {
		return errorBag
	}

	for i := 0; i < reflectType.NumField(); i++ {
		field := reflectType.Field(i)
		desc, inject := field.Tag.Lookup("query")

		if inject {
			name, required, targetType := parseInjectionParams(field, desc)
			valueField := fields.Field(i)

			prev := preprocess
			preprocess = func(writer http.ResponseWriter, r *http.Request, errorBag []error) []error {
				if r.URL.Query().Has(name) {
					parsedValue, err := parseValue(r.URL.Query().Get(name), targetType, field.Type)

					if err != nil {
						return append(errorBag, err)
					}

					valueField.Set(parsedValue)
				} else if required {
					return append(errorBag, errors.New("Missing required field "+name))
				}

				return prev(writer, r, errorBag)
			}
		}
	}

	return func(writer http.ResponseWriter, r *http.Request) {
		writer.Header().Set("Content-Type", "application/json")
		errorBag := preprocess(writer, r, nil)

		if len(errorBag) == 0 {
			handler.ServeHTTP(writer, r)
		} else {
			writer.WriteHeader(http.StatusBadRequest)
			body := struct {
				Errors []string `json:"errors"`
			}{
				Errors: make([]string, len(errorBag)),
			}

			for index, err := range errorBag {
				body.Errors[index] = err.Error()
			}
			payload, _ := json.Marshal(body)
			writer.Write(payload)
		}
	}
}

func parseInjectionParams(f reflect.StructField, params string) (string, bool, string) {
	parts := strings.Split(params, ",")
	parsedName := f.Name
	parsedType := f.Type.Name()
	required := false

	if len(parts) > 0 {
		parsedName = parts[0]
	}

	if len(parts) > 1 {
		required = parts[1] == "required"
	}

	if len(parts) > 2 {
		parsedType = parts[2]
	}

	return parsedName, required, parsedType
}

func parseValue(raw string, targetType string, fieldType reflect.Type) (reflect.Value, error) {
	if targetType[:2] == "[]" {
		return reflect.Zero(fieldType), errors.New("Multivalued fields are unsupported")
	}

	switch targetType {
	case "int":
		val, err := strconv.Atoi(raw)
		return reflect.ValueOf(val), err
	case "string":
		return reflect.ValueOf(raw), nil
	default:
		return reflect.Zero(fieldType), errors.New("Unsupported type: " + targetType)
	}
}
