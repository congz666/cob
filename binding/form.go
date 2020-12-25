package binding

import (
	"fmt"
	"net/http"
	"reflect"
	"strconv"
	"strings"
)

const defaultMemory = 32 << 20

type formPostBinding struct{}

func (formPostBinding) Bind(req *http.Request, obj interface{}) error {
	if err := req.ParseForm(); err != nil {
		return err
	}
	if err := req.ParseMultipartForm(defaultMemory); err != nil {
		if err != http.ErrNotMultipart {
			return err
		}
	}
	if err := process(req, obj); err != nil {
		return err
	}
	return nil
}

// process struct fields
func process(req *http.Request, obj interface{}) error {
	// Build map of fields keyed by effective name.
	fields := make(map[string]reflect.Value)
	v := reflect.ValueOf(obj).Elem() // the struct variable
	for i := 0; i < v.NumField(); i++ {
		fieldInfo := v.Type().Field(i) // a reflect.StructField
		tag := fieldInfo.Tag           // a reflect.StructTag
		name := tag.Get("form")
		if name == "" {
			name = strings.ToLower(fieldInfo.Name)
		}
		fields[name] = v.Field(i)
	}

	// Update struct field for each parameter in the request.
	for name, values := range req.Form {
		f := fields[name]
		if !f.IsValid() {
			continue // ignore unrecognized HTTP parameters
		}
		for _, value := range values {
			if f.Kind() == reflect.Slice {
				elem := reflect.New(f.Type().Elem()).Elem()
				if err := populate(elem, value); err != nil {
					return fmt.Errorf("%s: %v", name, err)
				}
				f.Set(reflect.Append(f, elem))
			} else {
				if err := populate(f, value); err != nil {
					return fmt.Errorf("%s: %v", name, err)
				}
			}
		}
	}
	return nil
}

// populate data in struct
func populate(v reflect.Value, value string) error {
	switch v.Kind() {
	case reflect.String:
		v.SetString(value)

	case reflect.Int, reflect.Int16, reflect.Int32, reflect.Int64:
		i, err := strconv.ParseInt(value, 10, 64)
		if err != nil {
			return err
		}
		v.SetInt(i)

	case reflect.Uint, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		i, err := strconv.ParseUint(value, 10, 64)
		if err != nil {
			return err
		}
		v.SetUint(i)

	case reflect.Float32, reflect.Float64:
		i, err := strconv.ParseFloat(value, 64)
		if err != nil {
			return err
		}
		v.SetFloat(i)

	case reflect.Bool:
		b, err := strconv.ParseBool(value)
		if err != nil {
			return err
		}
		v.SetBool(b)

	default:
		return fmt.Errorf("unsupported kind %s", v.Type())
	}
	return nil
}
