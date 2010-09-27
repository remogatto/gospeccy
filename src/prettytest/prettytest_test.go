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

func testAssertTrue(assert *Assertions) *Assertions {
	if assert.True(true).IsFailed() {
		assert.T.Errorf("True(true) should not fail\n")
	}
	if !assert.True(false).IsFailed() {
		assert.T.Errorf("True(false) should fail\n")
	}
	return assert
}

func testAssertFalse(assert *Assertions) *Assertions {
	if assert.False(false).IsFailed() {
		assert.T.Errorf("False(false) should not fail\n")
	}
	if !assert.False(true).IsFailed() {
		assert.T.Errorf("False(true) should fail\n")
	}
	return assert
}

func testAssertEqual(assert *Assertions) *Assertions {
	if assert.Equal("foo", "foo").IsFailed() {
		assert.T.Errorf("Equal(foo, foo) should not fail")
	}
	if !assert.Equal("foo", "bar").IsFailed() {
		assert.T.Errorf("Equal(foo, bar) should fail")
	}
	return assert
}

func TestBaseAssertions(t *testing.T) {
	DryRun(
		t,
		testAssertTrue,
		testAssertFalse,
		testAssertEqual,
	)
}

func testPending(assert *Assertions) *Assertions {
	return assert.Pending()
}

func testPass(assert *Assertions) *Assertions {
	assert.True(true)
	return assert
}

func TestRunner(t *testing.T) {
	Run(
		t,
		testPending,
		testPass,
	)
}

