package main

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestRunCmd(t *testing.T) {
	cmd := []string{"testdata/env", "/bin/bash", "testdata/echo.sh", "arg1=1", "arg2=2"}
	env := make(Environment)
	env["BAR"] = EnvValue{Value: "bar", NeedRemove: false}
	env["EMPTY"] = EnvValue{Value: "", NeedRemove: true}
	env["FOO"] = EnvValue{Value: "foo", NeedRemove: false}
	env["HELLO"] = EnvValue{Value: "hello", NeedRemove: false}
	env["UNSET"] = EnvValue{Value: "", NeedRemove: true}

	_ = os.Setenv("HELLO", "SHOULD_REPLACE")
	_ = os.Setenv("FOO", "SHOULD_REPLACE")
	_ = os.Setenv("UNSET", "SHOULD_REMOVE")
	_ = os.Setenv("ADDED", "from original env")
	_ = os.Setenv("EMPTY", "SHOULD_BE_EMPTY")

	os.Args = cmd

	assert.Equal(t, 0, RunCmd(cmd, env))
}
