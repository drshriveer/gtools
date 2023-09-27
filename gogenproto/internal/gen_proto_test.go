package internal

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGenProto(t *testing.T) {
	assert.FileExistsf(t, "./test.pb.go", "you must run go generate for this test to work.")
}
