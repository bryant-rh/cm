package main

import (
	"fmt"
	"os"

	"k8s.io/klog/v2"

	"github.com/bryant-rh/cm/cmd/client/cmd"
)

func main() {
	defer klog.Flush()
	cmd := cmd.NewCmd()
	if err := cmd.Execute(); err != nil {
		if klog.V(1).Enabled() {
			klog.Fatalf("%+v", err) // with stack trace
		} else {
			fmt.Fprintln(os.Stderr, err)
			os.Exit(1)
		}
	}
}
