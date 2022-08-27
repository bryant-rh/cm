package helmx

import (
	"fmt"
	"io/ioutil"
	"os"
	"reflect"
	"regexp"
	"strings"

	"github.com/go-courier/reflectx"
	"gopkg.in/yaml.v2"
)

func SetEnviron(envs map[string]string) {
	for k, v := range envs {
		if err := os.Setenv(k, v); err != nil {
			panic(err)
		}
	}
}

func LoadEnvConfig(envs map[string]string, environment string, feature string) {
	LoadEnvConfigFromFiles(envs, "master", feature)
	if environment != "" {
		LoadEnvConfigFromFiles(envs, environment, feature)
	}
}

func LoadEnvConfigFromFiles(envs map[string]string, envName string, feature string) map[string]string {
	loadEnvConfigFromFile(envs, envName)
	if feature != "" {
		loadEnvConfigFromFile(envs, envName+"-"+feature)
	}
	return envs
}

func loadEnvConfigFromFile(envs map[string]string, envName string) {
	filename := "config/" + strings.ToLower(envName) + ".yml"

	fmt.Fprintf(os.Stdout, "try to load env vars from %s ...\n", filename)

	envFileContent, err := ioutil.ReadFile(filename)
	if err == nil {
		var envVars map[string]string
		err := yaml.Unmarshal([]byte(envFileContent), &envVars)
		if err != nil {
			panic(err)
		}
		for key, value := range envVars {
			envs[key] = value
		}
	}
}

func MustMarshalToEnvs(envs map[string]string, v interface{}, prefix string) {
	rv := reflect.Indirect(reflect.ValueOf(v))
	tpe := rv.Type()

	for i := 0; i < tpe.NumField(); i++ {
		field := tpe.Field(i)
		env := field.Tag.Get("env")

		if len(env) > 0 {
			data, err := reflectx.MarshalText(rv.Field(i).Interface())
			if err != nil {
				continue
			}
			if len(data) > 0 {
				envs[strings.ToUpper(prefix+env)] = string(data)
			}
		}
	}
}

var reEnvVar = regexp.MustCompile(`(\$?\$)\{?([A-Za-z0-9_]+)\}?`)

func ResolveEnvVars(envVars map[string]string, s string) string {
	result := reEnvVar.ReplaceAllStringFunc(s, func(str string) string {
		matched := reEnvVar.FindAllStringSubmatch(str, -1)[0]

		// skip $${ }
		if matched[1] == "$$" {
			return "${" + matched[2] + "}"
		}

		if value, ok := envVars[matched[2]]; ok {
			return value
		}

		fmt.Fprintf(os.Stderr, "Missing environment variable ${%s}\n", matched[2])

		return "${" + matched[2] + "}"
	})

	return result
}
