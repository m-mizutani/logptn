package logptn_test

import (
	logptn "github.com/m-mizutani/logptn/lib"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestReadFile(t *testing.T) {
	g := logptn.Generator{}
	err := g.ReadFile("../test_data/t001.log")
	assert.Nil(t, err)
	logs := g.Logs()

	assert.Equal(t, 6, len(logs))
}

func TestReadFileWithMultipleNL(t *testing.T) {
	g := logptn.Generator{}
	err := g.ReadFile("../test_data/t002.log")
	assert.Nil(t, err)
	logs := g.Logs()

	assert.Equal(t, 6, len(logs))
}
