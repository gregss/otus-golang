package main

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestReadDir(t *testing.T) {
	t.Run("succses", func(t *testing.T) {
		expected := make(Environment)
		expected["BAR"] = EnvValue{Value: "bar", NeedRemove: false}
		expected["EMPTY"] = EnvValue{Value: "", NeedRemove: false}
		expected["FOO"] = EnvValue{Value: "   foo\nwith new line", NeedRemove: false}
		expected["HELLO"] = EnvValue{Value: "\"hello\"", NeedRemove: false}
		expected["UNSET"] = EnvValue{Value: "", NeedRemove: true}

		actual, err := ReadDir("testdata/env")

		assert.Nil(t, err)
		assert.Equal(t, expected, actual)
	})
}
