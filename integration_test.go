package main

import (
	"bytes"
	"io"
	"os"
	"regexp"
	"strings"
	"testing"
	"runtime"
)

func run(expr string) string {
	old := os.Stdout
	r, w, _ := os.Pipe()
	os.Stdout = w

	maxRoutines = runtime.NumCPU()
	execute_expression(expr)

	outc := make(chan string)
	go func() {
		var buf bytes.Buffer
		io.Copy(&buf, r)
		outc <- buf.String()
	}()

	w.Close()
	os.Stdout = old
	out := <-outc
	return strings.TrimSpace(unpad(out))
}

func expect(t *testing.T, output string, target string) {
	if target != output {
		t.Error("\nexpected:\n" + target + "\n\ngot:\n" + output)
	}
}

func unpad(result string) string {
	re := regexp.MustCompile("[ ]*\\n")
	return re.ReplaceAllString(result, "\n")
}

func TestEndswith(t *testing.T) {
	out := run("name in 'testdata/sub1' where name endswith '.txt'")
	expect(t, out, "test.txt\ntest2.txt")
}

func TestEndswithIgnorecase(t *testing.T) {
	out := run("name in 'testdata/sub1' where name ignorecase endswith '.TxT'")
	expect(t, out, "Test3.Txt\ntest.txt\ntest2.txt")
}

func TestStartswith(t *testing.T) {
	out := run("name in 'testdata/sub1' where name startswith 'test'")
	expect(t, out, "test.txt\ntest2.txt")
}

func TestStartswithIgnorecase(t *testing.T) {
	out := run("name in 'testdata/sub1' where name ignorecase startswith 'Test'")
	expect(t, out, "Test3.Txt\ntest.txt\ntest2.txt")
}

func TestIsfile(t *testing.T) {
	out := run("name in 'testdata/sub2' where isfile")
	expect(t, out, "a\nb\nc\nd")
}

func TestIsdir(t *testing.T) {
	out := run("name in 'testdata/sub2' where isdir")
	expect(t, out, "dir1/\ndir2/")
}

func TestNameEq(t *testing.T) {
	out := run("name in 'testdata/sub1' where name = 'test.txt'")
	expect(t, out, "test.txt")
}

func TestNameNeq(t *testing.T) {
	out := run("name in 'testdata/sub1' where name != 'test.txt'")
	expect(t, out, "Test3.Txt\ntest2.txt")
}

func TestNameContains(t *testing.T) {
	out := run("name in 'testdata/sub1' where name contains '2'")
	expect(t, out, "test2.txt")
}

func TestNameContainsIgnorecase(t *testing.T) {
	out := run("name in 'testdata/sub1' where name ignorecase contains 'Test'")
	expect(t, out, "Test3.Txt\ntest.txt\ntest2.txt")
}

func TestContentContains(t *testing.T) {
	out := run("name in 'testdata/sub1' where content contains 'some'")
	expect(t, out, "test.txt\ntest2.txt")
}

func TestContentContainsIgnorecase(t *testing.T) {
	out := run("name in 'testdata/sub1' where content ignorecase contains 'some'")
	expect(t, out, "Test3.Txt\ntest.txt\ntest2.txt")
}

func TestContentContainsIgnorecaseWithCaps(t *testing.T) {
	out := run("name in 'testdata/sub1' where content ignorecase contains 'SOME'")
	expect(t, out, "Test3.Txt\ntest.txt\ntest2.txt")
}

func TestContentEndswith(t *testing.T) {
	out := run("name in 'testdata/sub1' where content endswith '2'")
	expect(t, out, "test2.txt")
}

func TestContentEndswithWithContains(t *testing.T) {
	out := run("name in 'testdata/sub1' where content contains 'data' and content endswith '2'")
	expect(t, out, "test2.txt")
}

func TestContentStartsContainsAndEnds(t *testing.T) {
	out := run("name in 'testdata/sub1' where content startswith 'some' and content contains 'data' and content endswith '2'")
	expect(t, out, "test2.txt")
}

func TestPathExtraction(t *testing.T) {
	out := run("path in 'testdata/sub1' where name = 'test.txt'")
	if out != "testdata/sub1/test.txt" && out != "testdata\\sub1\\test.txt" {
		t.Error("\nexpected:\n" + "testdata/sub1/test.txt\n\ngot:\n" + out)
	}
}

func TestSizeExtraction(t *testing.T) {
	out := run("size in 'testdata/sub1' where name = 'test.txt'")
	expect(t, out, "10")
}

func TestCompoundSelect(t *testing.T) {
	out := run("name, size in 'testdata/sub1' where name = 'test.txt'")
	expect(t, out, "10          test.txt")
}

func TestOrExpression(t *testing.T) {
	out := run("name in 'testdata/sub2' where isfile or isdir")
	expect(t, out, "a\nb\ndir1/\nc\ndir2/\nd")
}

func TestAndExpression(t *testing.T) {
	out := run("name in 'testdata/sub1' where name startswith 'test.' and name endswith '.txt'")
	expect(t, out, "test.txt")
}

func TestNotExpression(t *testing.T) {
	out := run("name in 'testdata/sub1' where not name startswith 'test'")
	expect(t, out, "Test3.Txt")
}

func TestCompoundExpression(t *testing.T) {
	out := run("name in 'testdata/sub1' where name startswith 'T' or (name startswith 't' and not name contains '2')")
	expect(t, out, "Test3.Txt\ntest.txt")
}

func TestMultipleRoots(t *testing.T) {
	out := run("name in 'testdata/sub1', 'testdata/sub3' where name endswith '.txt'")
	expect(t, out, "test.txt\ntest2.txt\ntest4.txt")
}

func TestNameMatches(t *testing.T) {
	out := run("name in 'testdata/sub1' where name matches '[tT]est[0-9]'")
	expect(t, out, "Test3.Txt\ntest2.txt")
}

func TestPathMatches(t *testing.T) {
	out := run("name in 'testdata/sub1' where path matches '[tT]est[0-9]'")
	expect(t, out, "Test3.Txt\ntest2.txt")
}

func TestContentMatches(t *testing.T) {
	out := run("name in 'testdata/sub1' where content matches '[tT]ext'")
	expect(t, out, "Test3.Txt\ntest.txt")
}

func TestPathEqualsIgnorecase(t *testing.T) {
	out := run("name in 'testdata/sub1' where path ignorecase = 'testdata/sub1/test3.txt'")
	expect(t, out, "Test3.Txt")
}

func TestSha1Startswith(t *testing.T) {
	out := run("name in 'testdata' where sha1 startswith '3d47bc'")
	expect(t, out, "Test3.Txt")
}

func TestSha1Endswith(t *testing.T) {
	out := run("name in 'testdata' where sha1 endswith '8c308c'")
	expect(t, out, "Test3.Txt")
}

func TestSha1Contains(t *testing.T) {
	out := run("name in 'testdata' where sha1 contains 'efe0b9e47ab'")
	expect(t, out, "Test3.Txt")
}

func TestSha1Equals(t *testing.T) {
	out := run("name in 'testdata' where sha1 = '3d47bc8c8a81efe0b9e47ab4250f1a20ef8c308c'")
	expect(t, out, "Test3.Txt")
}

func TestSha1EqualsIgnorecase(t *testing.T) {
	out := run("name in 'testdata' where sha1 ignorecase = '3D47BC8c8a81efe0b9e47ab4250f1a20ef8c308c'")
	expect(t, out, "Test3.Txt")
}

func TestSha1Print(t *testing.T) {
	out := run("sha1 in 'testdata' where sha1 = '3d47bc8c8a81efe0b9e47ab4250f1a20ef8c308c'")
	expect(t, out, "3d47bc8c8a81efe0b9e47ab4250f1a20ef8c308c")
}

func TestMd5Startswith(t *testing.T) {
	out := run("name in 'testdata' where md5 startswith 'd1b0c3'")
	expect(t, out, "Test3.Txt")
}

func TestMd5Endswith(t *testing.T) {
	out := run("name in 'testdata' where md5 endswith '317d31'")
	expect(t, out, "Test3.Txt")
}

func TestMd5Contains(t *testing.T) {
	out := run("name in 'testdata' where md5 contains 'fd8d0f55'")
	expect(t, out, "Test3.Txt")
}

func TestMd5Equals(t *testing.T) {
	out := run("name in 'testdata' where md5 = 'd1b0c3ffb4dfd8d0f55a2a3d2a317d31'")
	expect(t, out, "Test3.Txt")
}

func TestMd5EqualsIgnorecase(t *testing.T) {
	out := run("name in 'testdata' where md5 ignorecase = 'D1B0C3Ffb4dfd8d0f55a2a3d2a317d31'")
	expect(t, out, "Test3.Txt")
}

func TestMd5Print(t *testing.T) {
	out := run("md5 in 'testdata' where md5 = 'd1b0c3ffb4dfd8d0f55a2a3d2a317d31'")
	expect(t, out, "d1b0c3ffb4dfd8d0f55a2a3d2a317d31")
}