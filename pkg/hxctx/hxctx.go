package hxctx

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"

	"github.com/davecgh/go-spew/spew"

	"github.com/logrusorgru/aurora"
)

type Context struct {
	Feature   string
	DeployEnv string
	Workspace string

	MainWorkspace *Workspace
	Workspaces    []*Workspace
}

func (ctx *Context) switchByRefName() {
	if commitRefName := os.Getenv("CI_COMMIT_REF_NAME"); commitRefName != "" {
		ctx.Feature, ctx.DeployEnv, ctx.Workspace = resolveFeatEnvWorkspace(commitRefName)
		spew.Dump(commitRefName, ctx.Feature, ctx.DeployEnv, ctx.Workspace)
	}
}

func (ctx *Context) setDefault(w *Workspace) {
	if v, ok := w.EnvVarSet.Values["PROJECT_NAME"]; ok && v != "" {
		w.Project.Name = v
	}

	if v, ok := w.EnvVarSet.Values["PROJECT_GROUP"]; ok && v != "" {
		w.Project.Group = v
	}

	if v, ok := w.EnvVarSet.Values["PROJECT_DESCRIPTION"]; ok && v != "" {
		w.Project.Description = v
	}

	if ctx.Feature != "" && w.Project.Feature == "" {
		w.Project.Feature = ctx.Feature
	}

	if ctx.DeployEnv != "" {
		w.DeployEnv = ctx.DeployEnv
	}
}

func (ctx *Context) Init() {
	ctx.switchByRefName()

	cwd, _ := os.Getwd()

	ls, err := ioutil.ReadDir(filepath.Join(cwd, "cmd"))
	if err != nil {
		w := NewWorkspace(cwd, cwd, ctx.setDefault)
		ctx.MainWorkspace = w
		ctx.Workspaces = append(ctx.Workspaces, w)
		return
	}

	for i := range ls {
		d := ls[i]

		if d.IsDir() {
			w := NewWorkspace(filepath.Join(cwd, "cmd", d.Name()), cwd, ctx.setDefault)
			if w.IsMain {
				ctx.MainWorkspace = w
			}
			ctx.Workspaces = append(ctx.Workspaces, w)
		}
	}

}

func (ctx *Context) Range(fn func(w *Workspace)) {
	do := func(w *Workspace) {
		fmt.Fprintln(os.Stdout, aurora.Sprintf(aurora.Magenta(`
===== workspace %s (%s) ======

`), w.Name, w.WorkspaceRoot))

		fn(w)
	}

	for i := range ctx.Workspaces {
		workspace := ctx.Workspaces[i]

		if ctx.Workspace != "" && workspace.Name != ctx.Workspace {
			continue
		}

		if workspace.IsMain {
			do(workspace)
			break
		}
	}

	for i := range ctx.Workspaces {
		workspace := ctx.Workspaces[i]

		if ctx.Workspace != "" && workspace.Name != ctx.Workspace && !strings.HasSuffix(workspace.Name, ctx.Workspace) {
			continue
		}

		if workspace.IsMain {
			continue
		}

		do(workspace)
	}
}
