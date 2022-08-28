package cmd

import (
	"github.com/bryant-rh/cm/cmd/client/global"
	"github.com/bryant-rh/cm/pkg/client"
	"github.com/bryant-rh/cm/pkg/environment"
	"github.com/bryant-rh/cm/pkg/util"
	"flag"
	"fmt"
	"math/rand"
	"os"
	"strings"
	"time"

	"github.com/pkg/errors"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
	"github.com/spf13/viper"
	"k8s.io/klog/v2"
)



func NewCmd() *cobra.Command {
	// rootCmd represents the base command when called without any subcommands
	rootCmd := &cobra.Command{
		Use:     "cm", // This is prefixed by kubectl in the custom usage template
		Short:   "cm is a tool for managing k8s clusters",
		Long: `cm is a tool for managing k8s clusters.
You can invoke cm through comand: "cm [command]..."`,
		SilenceUsage:  true,
		SilenceErrors: true,
		//PersistentPreRunE: preRun,
		//PersistentPostRun: showUpgradeNotification,
		CompletionOptions: cobra.CompletionOptions{
			DisableDefaultCmd: true,
		},
	}
	global.CMClient = client.NewGithubClient(global.CM_SERVER_BASEURL)
	global.ProxyClient = client.NewGithubClient(global.KUBE_TUNNEL_GATEWAY_HOST)

	rootCmd.AddCommand(NewCmdVersion())
	rootCmd.AddCommand(NewCmdGet())
	rootCmd.AddCommand(NewCmdCreate())
	rootCmd.AddCommand(NewCmdUpdate())
	rootCmd.AddCommand(NewCmdDelete())
	rootCmd.AddCommand(NewCmdAdd())
	rootCmd.AddCommand(NewCmdHx())

	return rootCmd
}

func initLog() {
	klog.InitFlags(nil)
	rand.Seed(time.Now().UnixNano())

	pflag.CommandLine.AddGoFlagSet(flag.CommandLine)
	_ = flag.CommandLine.Parse([]string{}) // convince pkg/flag we parsed the flags

	flag.CommandLine.VisitAll(func(f *flag.Flag) {
		if f.Name != "v" { // hide all glog flags except for -v
			pflag.Lookup(f.Name).Hidden = true
		}
	})
	if err := flag.Set("logtostderr", "true"); err != nil {
		fmt.Printf("can't set log to stderr %+v", err)
		os.Exit(1)
	}
}

func initConfig() {

	global.Paths = environment.MustGetCmPaths()
	if !util.Exists(global.Paths.BasePath()) {
		if err := ensureDirs(global.Paths.BasePath()); err != nil {
			klog.Fatal(err)
		}
	}

	cm_config := fmt.Sprintf("%s/cm.yaml", global.Paths.BasePath())
	if !util.Exists(cm_config) {
		f, err := os.Create(cm_config)
		if err != nil {
			klog.Fatal(err)
		}
		f.Close()
	}
	viper.AddConfigPath(global.Paths.BasePath())
	viper.SetConfigName("cm")
	viper.SetConfigType("yaml")
	viper.AutomaticEnv()
	replacer := strings.NewReplacer(".", "_")
	viper.SetEnvKeyReplacer(replacer)
	if err := viper.ReadInConfig(); err != nil {
		klog.Fatalf("read config file failed! %v;\n 请在目录:[%s] 创建配置文件: [cm.yaml], 配置CM_SERVER_BASEURL、CM_SERVER_USERNAME、CM_SERVER_PASSWORD 或者配置对应环境变量\n", err, global.Paths.BasePath())

	}

	global.CM_SERVER_BASEURL = viper.GetString("CM_SERVER_BASEURL")
	global.CM_SERVER_USERNAME = viper.GetString("CM_SERVER_USERNAME")
	global.CM_SERVER_PASSWORD = viper.GetString("CM_SERVER_PASSWORD")

	global.KUBE_TUNNEL_GATEWAY_HOST = viper.GetString("KUBE_TUNNEL_GATEWAY_HOST")

}

func init() {
	initLog()
	initConfig()
}

func ensureDirs(paths ...string) error {
	for _, p := range paths {
		klog.V(4).Infof("Ensure creating dir: %q", p)
		if err := os.MkdirAll(p, 0o755); err != nil {
			return errors.Wrapf(err, "failed to ensure create directory %q", p)
		}
	}
	return nil
}
