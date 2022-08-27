package helmx

import (
	"io/ioutil"
	"os"
	"strings"

	"gopkg.in/yaml.v2"
)

var (
	EnvVarKeyCommitSha        = "CI_COMMIT_SHA"
	EnvVarKeyCommitRefName    = "CI_COMMIT_REF_NAME"
	EnvVarKeyProjectNamespace = "CI_PROJECT_NAMESPACE"
	EnvVarKeyProjectName      = "CI_PROJECT_NAME"
	EnvVarKeyCommitMessage    = "CI_COMMIT_MESSAGE"
)

func GetFeatAndEnv(commitRef string) (feat string, env string) {
	if commitRef == "develop" {
		return "", "staging"
	}

	if commitRef == "master" {
		return "", ""
	}

	parts := strings.Split(commitRef, ".")

	featParts := strings.Split(parts[0], "/")

	if len(featParts) != 2 {
		return
	}

	if strings.HasPrefix(featParts[0], "feat") {
		feat = strings.ToLower(featParts[1])
		env = "staging"

		if len(parts) == 2 {
			env = strings.ToLower(parts[1])
			if env == "demo" {
				env = ""
			}
		}
	}

	return
}

func GetShortSha() string {
	ref := os.Getenv(EnvVarKeyCommitSha)
	if len(ref) >= 7 {
		return ref[0:7]
	}
	return ""
}

func MustWriteYAML(filename string, v interface{}) {
	data, err := yaml.Marshal(v)
	if err != nil {
		panic(err)
	}
	err = ioutil.WriteFile(filename, data, os.ModePerm)
	if err != nil {
		panic(err)
	}
}

// func ptrInt32(i int32) *int32 {
// 	return &i
// }
func isPathExist(path string) bool {
	f, _ := os.Stat(path)
	return f != nil
}
