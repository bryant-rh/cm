package hxctx

import (
	"io/ioutil"
	"os"
	"strings"

	"gopkg.in/yaml.v2"
)

func resolveFeatEnvWorkspace(commitRef string) (feat string, env string, workspace string) {
	if commitRef == "develop" || commitRef == "main" {
		return "", "staging", ""
	}

	if commitRef == "master" {
		return "", "", ""
	}

	defer func() {
		if env == "demo" {
			env = ""
		}
	}()

	parts := strings.Split(commitRef, ".")

	if len(parts) == 2 {
		env = strings.ToLower(parts[1])
	}

	if env == "" {
		env = "staging"
	}

	metaParts := strings.Split(parts[0], "/")

	if len(metaParts)%2 != 0 {
		return
	}

	for i := 0; i < len(metaParts)/2; i++ {
		t := metaParts[2*i]

		if t == "release" {
			env = "demo"
			feat = strings.ToLower(metaParts[2*i+1])
		}

		if t == "feat" || t == "feature" || t == "f" ||
			t == "test" || t == "t" {
			feat = strings.ToLower(metaParts[2*i+1])
		}

		if t == "workspace" || t == "w" {
			workspace = strings.ToLower(metaParts[2*i+1])
		}
	}

	if feat == "" {
		nameFeat := strings.SplitN(workspace, "--", 2)

		if len(nameFeat) == 2 {
			workspace = nameFeat[0]
			feat = nameFeat[1]
		}
	}

	return
}

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
