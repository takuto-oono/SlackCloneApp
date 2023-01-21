package test

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSomething(t *testing.T) {
	assert := assert.New(t)

	var a string = "Hello Golang"
	var b string = "Hello Golang"

	assert.Equal(a, b, "The two words should be the same.")
}
