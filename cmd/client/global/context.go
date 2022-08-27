package global

import (
	"github.com/go-courier/helmx/spec"
)

type Context struct {
	spec.Spec
	DefaultEnvs      map[string]string
	Environment      string
	Version          string
	Templates        map[string][]byte
	ReferFiles       Refers
	CommitShaMessage string
	CommitSha        string
}
