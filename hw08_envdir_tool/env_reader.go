package main

import (
	"bufio"
	"bytes"
	"io/ioutil"
	"os"
	"strings"
)

type Environment map[string]EnvValue

// EnvValue helps to distinguish between empty files and files with the first empty line.
type EnvValue struct {
	Value      string
	NeedRemove bool
}

// ReadDir reads a specified directory and returns map of env variables.
// Variables represented as files where filename is name of variable, file first line is a value.
func ReadDir(dir string) (Environment, error) {
	// Place your code here
	files, err := ioutil.ReadDir(dir)
	if err != nil {
		return nil, err
	}

	env := make(Environment)

	for _, fstat := range files {
		fname := fstat.Name()

		file, err := os.Open(dir + "/" + fname)
		if err != nil {
			return nil, err
		}
		defer file.Close()

		scanner := bufio.NewScanner(file)
		scanner.Scan()
		b := scanner.Bytes()
		b = bytes.ReplaceAll(b, []byte("\x00"), []byte("\n"))
		s := string(b)

		if s == "" {
			env[fname] = EnvValue{Value: s, NeedRemove: true}
			continue
		} else {
			if strings.TrimSpace(s) == "" {
				s = ""
			}
			env[fname] = EnvValue{
				Value:      s,
				NeedRemove: false,
			}
		}
	}

	return env, nil
}
