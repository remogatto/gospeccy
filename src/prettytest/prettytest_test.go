package prettytest

import (
	"testing"
)

func testAssertTrue(assert *assertions) {
	if !assert.True(true) {
		assert.T.Errorf("True(true) should return true\n")
	}
}

func testAssertFalse(assert *assertions) {
	if !assert.False(false) {
		assert.T.Errorf("False(false) should return true\n")
	}
}

func testAssertEqual(assert *assertions) {
	assert.True(assert.Equal("foo", "foo"))
	assert.False(assert.Equal("foo", "bar"))
}

func TestPrettyTest(t *testing.T) {
	Run(
		t,
		"true/false Assertions",
		testAssertTrue,
		testAssertFalse,
	)

	DryRun(
		t,
		testAssertEqual,
	)
		
}
