package hxctx

import (
	"reflect"
	"strings"

	encodingx "github.com/go-courier/x/encoding"
)

func MustMarshalToEnvs(envs map[string]string, v interface{}, prefix string) {
	rv := reflect.Indirect(reflect.ValueOf(v))

	var walk func(rv reflect.Value)

	walk = func(rv reflect.Value) {
		rv = reflect.Indirect(rv)
		tpe := rv.Type()

		for i := 0; i < tpe.NumField(); i++ {
			field := tpe.Field(i)
			fieldValue := rv.Field(i)

			env, ok := field.Tag.Lookup("env")

			if ok {
				if env == "" {
					env = field.Name
				}

				data, err := encodingx.MarshalText(fieldValue.Interface())
				if err != nil {
					continue
				}

				if len(data) > 0 {
					envs[strings.ToUpper(prefix+env)] = string(data)
				}
			} else if field.Anonymous {
				walk(fieldValue)
				continue
			}
		}
	}

	walk(rv)
}
