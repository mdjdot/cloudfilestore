package util

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSha1Stream_Update(t *testing.T) {
	data := []byte("asat")
	ss := &Sha1Stream{}
	ss.Update(data)
	assert.NotNil(t, ss._sha1)
}
