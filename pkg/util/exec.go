package util

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/alecthomas/chroma/quick"
	"mvdan.cc/sh/v3/interp"
	"mvdan.cc/sh/v3/syntax"
)

func StdRun(args ...string) {
	script := strings.Join(args, " ")

	fmt.Printf("$ %s\n", script)

	file, _ := syntax.NewParser().Parse(strings.NewReader(script), "")
	runner, err := interp.New(
		interp.StdIO(os.Stdin, os.Stdout, os.Stderr),
	)
	if err != nil {
		panic(err)
	}
	if err := runner.Run(context.Background(), file); err != nil {
		panic(err)
	}
}

func PrintifyJSON(writer io.Writer, v interface{}) error {
	data, _ := json.MarshalIndent(v, "", "  ")
	_, _ = io.WriteString(writer, "\n```\n")
	defer func() {
		_, _ = io.WriteString(writer, "\n```\n")
	}()

	if err := quick.Highlight(writer, string(data), "json", "terminal", "vim"); err != nil {
		panic(err)
	}
	return nil
}

func PrintifyDockerfile(writer io.Writer, dockerfile []byte) {
	_, _ = io.WriteString(writer, "\n```\n")
	defer func() {
		_, _ = io.WriteString(writer, "\n```\n")
	}()

	if err := quick.Highlight(writer, string(dockerfile), "Docker", "terminal", "vim"); err != nil {
		panic(err)
	}
}
