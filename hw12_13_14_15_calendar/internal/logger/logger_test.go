package logger

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestLogger(t *testing.T) {
	fileName := "filename"
	(New("info", fileName)).Info("test")
	defer os.Remove(fileName)
	b, _ := os.ReadFile(fileName)
	s := string(b)
	assert.Equal(t, "test\n", s)
}
