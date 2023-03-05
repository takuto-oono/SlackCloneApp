package models

import (
	"math/rand"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/xyproto/randomstring"
)

func TestCreateDirectMessage(t *testing.T) {
	// if testing.Short() {
	// 	t.Skip("skipping test in short mode.")
	// }
	
	dm := NewDirectMessage(randomstring.EnglishFrequencyString(99), rand.Uint32(), uint(rand.Uint64()))
	res := dm.Create()
	assert.NotEqual(t, 0, dm.ID)
	assert.Empty(t, res.Error)
}
