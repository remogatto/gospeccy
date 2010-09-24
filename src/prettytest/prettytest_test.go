/*

Copyright (c) 2010 Andrea Fazzi

Permission is hereby granted, free of charge, to any person obtaining
a copy of this software and associated documentation files (the
"Software"), to deal in the Software without restriction, including
without limitation the rights to use, copy, modify, merge, publish,
distribute, sublicense, and/or sell copies of the Software, and to
permit persons to whom the Software is furnished to do so, subject to
the following conditions:

The above copyright notice and this permission notice shall be
included in all copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND,
EXPRESS OR IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF
MERCHANTABILITY, FITNESS FOR A PARTICULAR PURPOSE AND
NONINFRINGEMENT. IN NO EVENT SHALL THE AUTHORS OR COPYRIGHT HOLDERS BE
LIABLE FOR ANY CLAIM, DAMAGES OR OTHER LIABILITY, WHETHER IN AN ACTION
OF CONTRACT, TORT OR OTHERWISE, ARISING FROM, OUT OF OR IN CONNECTION
WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE SOFTWARE.

*/

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

func testPending(assert *assertions) {
	assert.Pending("testPending")
}

func TestPrettyTest(t *testing.T) {
	Run(
		t,
		"TestPrettyTest",
		testAssertTrue,
		testAssertFalse,
		testPending,
	)

	DryRun(
		t,
		testAssertEqual,
	)
		
}
