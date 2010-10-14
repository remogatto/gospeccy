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
	"os"
	"io/ioutil"
)

func testAssertTrue(assert *T) {
	assert.True(true)
	if assert.Failed() {
		assert.T.Errorf("True(true) should not fail\n")
	}

	assert.True(false)
	if !assert.Failed() {
		assert.T.Errorf("True(false) should fail\n")
	}
}

func testAssertFalse(assert *T) {
	assert.False(false)
	if assert.Failed() {
		assert.T.Errorf("False(false) should not fail\n")
	}

	assert.False(true)
	if !assert.Failed() {
		assert.T.Errorf("False(true) should fail\n")
	}
}

func testAssertEqual(assert *T) {
	assert.Equal("foo", "foo")
	if assert.Failed() {
		assert.T.Errorf("Equal(foo, foo) should not fail")
	}

	assert.Equal("foo", "bar")
	if !assert.Failed() {
		assert.T.Errorf("Equal(foo, bar) should fail")
	}
}

func testLastAssertionStatus(assert *T) {
	assert.Equal("foo", "bar")
	assert.Equal("foo", "foo")

	if assert.Failed() {
		assert.T.Errorf("Assertion last status should not be STATUS_FAIL")
	}

	if !assert.TestFailed() {
		assert.T.Errorf("Test status should be STATUS_FAIL")
	}

}

func TestBaseAssertions(t *testing.T) {
	DryRun(
		t,
		testAssertTrue,
		testAssertFalse,
		testAssertEqual,
		testLastAssertionStatus,
	)
}

func testPending(assert *T) { 
	assert.Pending()
}

func testPass(assert *T) { 
	assert.True(true)
}

var state int = 0

func before(assert *T) {
	state += 2
}

func after(assert *T) {
	state--
}

func beforeAll(assert *T) {
	state = 0
}

func afterAll(assert *T) {
	state = 0
}

func testSetup_1(assert *T) {
	assert.Equal(2, state)
}

func testSetup_2(assert *T) {
	assert.Equal(3, state)
}

func TestRunner(t *testing.T) {
	Run(
		t,
		testPending,
		testPass,
	)
}

func TestSetupTeardown(t *testing.T) {
	Run(
		t,
		before,
		after,
		testSetup_1,
		testSetup_2,
	)
}

func TestMisplacedSetupTeardown(t *testing.T) {
	state = 0
	Run(
		t,
		testSetup_1,
		before,
		testSetup_2,
		after,
	)	
}

func TestSetupAllTeardownAll(t *testing.T) {
	state = 10
	Run(
		t,
		beforeAll,
		afterAll,
		before,
		after,
		testSetup_1,
		testSetup_2,
	)
	if state != 0 {
		t.Errorf("state should be 0 afterAll tests\n")
	}
}

func afterTestPath(assert *T) {
	os.Remove("testfile")
}

func testPath(assert *T) {
	ioutil.WriteFile("testfile", nil, 0600)
	assert.Path("testfile")
	
	assert.Dry = true
	assert.Path("foo")
	assert.True(assert.Failed())
}

func TestPath(t *testing.T) {
	Run(
		t,
		afterTestPath,
		testPath,
	)
}

func testNil(assert *T) { assert.Nil(nil) }

func TestNil(t *testing.T) {
	Run(
		t,
		testNil,
	)
}
