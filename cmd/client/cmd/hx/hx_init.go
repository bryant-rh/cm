package hx

import (
	"fmt"
	"os"
	"path"
	"strings"

	"github.com/bryant-rh/cm/cmd/client/global"
	"github.com/bryant-rh/cm/pkg/helmx"
	"github.com/bryant-rh/cm/pkg/temp"
	"github.com/bryant-rh/cm/pkg/util"

	"github.com/go-courier/helmx/spec"
	"github.com/spf13/cobra"
	"gopkg.in/AlecAivazis/survey.v1"
	"k8s.io/klog/v2"
)

func NewCmdHxInit() *cobra.Command {
	// certCmd represents the cert command
	cmd := &cobra.Command{
		Use:   "init",
		Short: "project initial",
		Run: func(cmd *cobra.Command, args []string) {
			klog.V(4).Infoln("Hx Init")

			cwd, _ := os.Getwd()

			if len(global.Ctx.Project.Name) == 0 {
				global.Ctx.Project.Name = path.Base(cwd)
			}

			if len(global.Ctx.Project.Group) == 0 {
				parts := strings.Split(path.Dir(cwd), "/")
				for i, part := range parts {
					if strings.HasPrefix(part, "git") && i < len(parts)-1 {
						global.Ctx.Project.Group = parts[i+1]
					}
				}

			}

			if len(global.Ctx.Project.Description) == 0 {
				global.Ctx.Project.Description = global.Ctx.Project.Name
			}

			answers := struct {
				Name        string
				Group       string
				Description string
				Version     string
			}{}

			questings := []*survey.Question{
				{
					Name: "name",
					Prompt: &survey.Input{
						Message: "项目名称",
						Default: global.Ctx.Project.Name,
					},
					Validate: survey.Required,
				},
				{
					Name: "description",
					Prompt: &survey.Input{
						Message: "项目描述",
						Default: global.Ctx.Project.Description,
					},
					Validate: survey.Required,
				},
				{
					Name: "group",
					Prompt: &survey.Input{
						Message: "项目所属部门",
						Default: global.Ctx.Project.Group,
					},
					Validate: survey.Required,
				},
				{
					Name: "version",
					Prompt: &survey.Input{
						Message: "项目版本号 (x.x.x)",
						Default: global.Ctx.Project.Version.String(),
					},
					Validate: survey.Required,
				},
			}

			err := survey.Ask(questings, &answers)
			if err != nil {
				fmt.Println(err.Error())
				return
			}

			global.Ctx.Project.Name = answers.Name
			global.Ctx.Project.Group = answers.Group
			global.Ctx.Project.Description = answers.Description

			if answers.Version != "" {
				version, err := spec.ParseVersion(answers.Version)
				if err == nil {
					global.Ctx.Project.Version = *version
				}
			}

			helmx.MustWriteYAML(global.HelmxRootFile, struct {
				Project *spec.Project
			}{
				Project: global.Ctx.Spec.Project,
			})

			//生成dockerfile.default
			file_type := ""
			switch {
			case util.Exists("./go.mod"):
				file_type = "go"
			case util.Exists("./pom.xml"):
				file_type = "java"
			case util.Exists("./requirements.txt"):
				file_type = "python"
			}
			fmt.Println(util.GreenColor("生成helmx.yml"))
			_ = util.WriteToFile("./helmx.yml", temp.Helmxfile(), func(v interface{}) ([]byte, error) {
				return v.([]byte), nil
			})

			fmt.Println(util.GreenColor("生成dockerfile.default"))
			_ = util.WriteToFile("./Dockerfile.default", temp.Dockerfile(file_type, global.Ctx.Project.Name), func(v interface{}) ([]byte, error) {
				return v.([]byte), nil
			})

			if file_type == "go" {
				fmt.Println(util.GreenColor("生成Makefile.default"))
				_ = util.WriteToFile("./Makefile.default", temp.Makefile(), func(v interface{}) ([]byte, error) {
					return v.([]byte), nil
				})

			}

			fmt.Println(util.GreenColor("项目初始化完成"))
		},
	}

	return cmd
}
