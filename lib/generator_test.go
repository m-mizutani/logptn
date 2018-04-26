package logptn_test

import (
	logptn "github.com/m-mizutani/logptn/lib"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestReadFile(t *testing.T) {
	g := logptn.NewGenerator()
	err := g.ReadFile("../test_data/t001.log")
	assert.Nil(t, err)
	logs := g.Logs()

	assert.Equal(t, 6, len(logs))
	assert.Equal(t, "9", logs[0].Chunk[3].Data)
}

func TestReadFileWithMultipleNL(t *testing.T) {
	g := logptn.NewGenerator()
	err := g.ReadFile("../test_data/t002.log")
	assert.Nil(t, err)
	logs := g.Logs()

	assert.Equal(t, 6, len(logs))
}

func TestFinalize(t *testing.T) {
	g := logptn.NewGenerator()
	err := g.ReadFile("../test_data/t001.log")
	assert.Nil(t, err)
	g.Finalize()

	formats := g.Formats()
	assert.True(t, len(formats) > 0)
}
