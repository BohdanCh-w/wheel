package assert_test

import (
	"testing"

	"github.com/bohdanch-w/wheel/assert"
)

func TestAssert(t *testing.T) {
	assert.Assert(true, "this should not fail")

	assert.Assert(false, "this should fail")
}
