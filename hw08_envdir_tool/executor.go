package main

import (
	"os"
	"os/exec"
	"strings"

	"github.com/spf13/cobra"
)

// RunCmd runs a command + arguments (cmd) with environment variables from env.
func RunCmd(cmd []string, env Environment) (returnCode int) {
	returnCode = 0
	rootCmd := &cobra.Command{
		Use: cmd[0],
		Run: func(cmd *cobra.Command, args []string) {
			returnCode = run(args, env)
		},
	}

	if err := rootCmd.Execute(); err != nil {
		returnCode = 1
		return
	}

	return
}

func run(args []string, env Environment) (returnCode int) {
	returnCode = 0

	c := exec.Command(args[1], args[2:]...) //#nosec G204

	envM := make(map[string]string)
	for _, es := range os.Environ() {
		e := strings.SplitN(es, "=", 2)
		envM[e[0]] = es
	}

	for k, v := range env {
		if _, ok := envM[k]; ok {
			if !v.NeedRemove {
				envM[k] = k + "=" + v.Value
			} else {
				delete(envM, k)
			}
		} else {
			envM[k] = k + "=" + v.Value
		}
	}

	slenv := make([]string, 0, len(envM))
	for _, v := range envM {
		slenv = append(slenv, v)
	}

	c.Env = slenv
	c.Stdout = os.Stdout
	c.Stderr = os.Stderr
	if err := c.Run(); err != nil {
		returnCode = 1
		return
	}

	return
}
