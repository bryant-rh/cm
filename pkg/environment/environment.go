package environment

import (
	"os"
	"path/filepath"

	"github.com/pkg/errors"
	"k8s.io/client-go/util/homedir"
	"k8s.io/klog/v2"
)

// Paths contains all important environment paths
type Paths struct {
	base string
}

// MustGetCmPaths returns the inferred paths for krew. By default, it assumes
// $HOME/.cm as the base path, but can be overridden via CM_ROOT environment
// variable.
func MustGetCmPaths() Paths {
	base := filepath.Join(homedir.HomeDir(), ".cm")
	if fromEnv := os.Getenv("CM_ROOT"); fromEnv != "" {
		base = fromEnv
		klog.V(4).Infof("using environment override CM_ROOT=%s", fromEnv)
	}
	base, err := filepath.Abs(base)
	if err != nil {
		panic(errors.Wrap(err, "cannot get absolute path"))
	}
	return NewPaths(base)
}

func NewPaths(base string) Paths {
	return Paths{base: base}
}

// BasePath returns krew base directory.
func (p Paths) BasePath() string { return p.base }
