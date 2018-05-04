package logptn_test

import (
	logptn "github.com/m-mizutani/logptn/lib"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestSplit1(t *testing.T) {
	s := logptn.NewSplitter()
	c := s.Split("a [b]")
	assert.Equal(t, 5, len(c))
	assert.Equal(t, "b", c[3].Data)
	assert.Equal(t, "]", c[4].Data)
}

func TestSplit2(t *testing.T) {
	s := logptn.NewSplitter()
	c := s.Split("abc [bdw]")
	assert.Equal(t, 5, len(c))
	assert.Equal(t, "bdw", c[3].Data)
	assert.Equal(t, "]", c[4].Data)
}

func TestSplit3(t *testing.T) {
	s := logptn.NewSplitter()
	c := s.Split(" abc [bdw]")
	assert.Equal(t, 6, len(c))
	assert.Equal(t, " ", c[0].Data)
	assert.Equal(t, "bdw", c[4].Data)
	assert.Equal(t, "]", c[5].Data)
}

func TestSplitWithDelim(t *testing.T) {
	s := logptn.NewSplitter()
	s.SetDelim("XYZ")
	c := s.Split(" abc [Xbdw]")
	assert.Equal(t, 3, len(c))
	assert.Equal(t, " abc [", c[0].Data)
	assert.Equal(t, "X", c[1].Data)
	assert.Equal(t, "bdw]", c[2].Data)
}

func TestHeuristicsSplit1(t *testing.T) {
	s := logptn.NewSplitter()
	c := s.Split("xx hoge<hoge@example.com>")
	assert.Equal(t, 6, len(c))
	assert.Equal(t, "hoge@example.com", c[4].Data)
}

func TestHeuristicsSplit2(t *testing.T) {
	s := logptn.NewSplitter()
	c := s.Split("2018-05-05 12:23:34 Yes")
	assert.Equal(t, 5, len(c))
	assert.Equal(t, "12:23:34", c[2].Data)
}

func TestHeuristicsSplit3(t *testing.T) {
	s := logptn.NewSplitter()
	c := s.Split("2018-05-05T12:23:34 Yes")
	assert.Equal(t, 3, len(c))
	assert.Equal(t, "Yes", c[2].Data)
}
