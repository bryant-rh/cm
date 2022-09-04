package hxctx

import (
	"bytes"
	"fmt"
	"os"
	"regexp"
	"strconv"
)

var reEnvVar = regexp.MustCompile(`(\$?\$)\{?([A-Za-z0-9_]+)\}?`)

func ResolveEnvVars(envVars map[string]string, data []byte) []byte {
	result := reEnvVar.ReplaceAllFunc(data, func(str []byte) []byte {
		matched := reEnvVar.FindAllSubmatch(str, -1)[0]
		key := string(matched[2])

		// skip $${ }
		if string(matched[1]) == "$$" {
			return []byte("${" + key + "}")
		}

		if value, ok := envVars[key]; ok {
			return []byte(value)
		}

		fmt.Fprintf(os.Stderr, "Missing environment variable ${%s}\n", key)

		return []byte("${" + key + "}")
	})

	return result
}

func SetEnviron(envs map[string]string) {
	for k, v := range envs {
		if err := os.Setenv(k, v); err != nil {
			panic(err)
		}
	}
}

var reVariables = regexp.MustCompile(`\$\{\{([^}]+)\}\}`)

func NewTemplate(tpl []byte) (*Template, error) {
	return &Template{tpl: tpl}, nil
}

type Template struct {
	tpl []byte
}

func (t *Template) Execute(values map[string]string) []byte {
	return reVariables.ReplaceAllFunc(t.tpl, func(str []byte) []byte {
		matched := reVariables.FindAllSubmatch(str, -1)[0]

		key := string(bytes.TrimSpace(matched[1]))

		if value, ok := values[key]; ok {
			return []byte(strconv.Quote(value))
		}

		return []byte("${{ " + key + " }}")
	})
}

func MustNewTemplate(data []byte) *Template {
	t, err := NewTemplate(data)
	if err != nil {
		panic(err)
	}
	return t
}
