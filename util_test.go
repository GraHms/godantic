package godantic

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func Test_JoinKeys(t *testing.T) {
	res := joinKeys("parent", "child")
	assert.Equal(t, res, "parent.child")
}
