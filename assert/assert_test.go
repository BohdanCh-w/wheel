package assert_test

import (
	"testing"

	"github.com/bohdanch-w/wheel/assert"
)

func TestAssert(_ *testing.T) {
	assert.Assert(true, "this should not fail")

	assert.Assert(false, "this should fail")
}
