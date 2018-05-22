package main

/*
	Copyright (c) 2009 The Go Authors. All rights reserved.

	Redistribution and use in source and binary forms, with or without
	modification, are permitted provided that the following conditions are
	met:

	   * Redistributions of source code must retain the above copyright
	notice, this list of conditions and the following disclaimer.
	   * Redistributions in binary form must reproduce the above
	copyright notice, this list of conditions and the following disclaimer
	in the documentation and/or other materials provided with the
	distribution.
	   * Neither the name of Google Inc. nor the names of its
	contributors may be used to endorse or promote products derived from
	this software without specific prior written permission.

	THIS SOFTWARE IS PROVIDED BY THE COPYRIGHT HOLDERS AND CONTRIBUTORS
	"AS IS" AND ANY EXPRESS OR IMPLIED WARRANTIES, INCLUDING, BUT NOT
	LIMITED TO, THE IMPLIED WARRANTIES OF MERCHANTABILITY AND FITNESS FOR
	A PARTICULAR PURPOSE ARE DISCLAIMED. IN NO EVENT SHALL THE COPYRIGHT
	OWNER OR CONTRIBUTORS BE LIABLE FOR ANY DIRECT, INDIRECT, INCIDENTAL,
	SPECIAL, EXEMPLARY, OR CONSEQUENTIAL DAMAGES (INCLUDING, BUT NOT
	LIMITED TO, PROCUREMENT OF SUBSTITUTE GOODS OR SERVICES; LOSS OF USE,
	DATA, OR PROFITS; OR BUSINESS INTERRUPTION) HOWEVER CAUSED AND ON ANY
	THEORY OF LIABILITY, WHETHER IN CONTRACT, STRICT LIABILITY, OR TORT
	(INCLUDING NEGLIGENCE OR OTHERWISE) ARISING IN ANY WAY OUT OF THE USE
	OF THIS SOFTWARE, EVEN IF ADVISED OF THE POSSIBILITY OF SUCH DAMAGE.
*/

import (
	"bytes"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
	"testing"
)

// code stolen from https://golang.org/src/cmd/gofmt/gofmt_test.go


func TestBackupFile(t *testing.T) {
	dir, err := ioutil.TempDir("", "gofmt_test")
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(dir)
	name, err := backupFile(filepath.Join(dir, "foo.go"), []byte("  package main"), 0644)
	if err != nil {
		t.Fatal(err)
	}
	t.Logf("Created: %s", name)
}

func TestDiff(t *testing.T) {
	if _, err := exec.LookPath("diff"); err != nil {
		t.Skipf("skip test on %s: diff command is required", runtime.GOOS)
	}
	in := []byte("first\nsecond\n")
	out := []byte("first\nthird\n")
	filename := "difftest.txt"
	b, err := diff(in, out, filename)
	if err != nil {
		t.Fatal(err)
	}

	if runtime.GOOS == "windows" {
		b = bytes.ReplaceAll(b, []byte{'\r', '\n'}, []byte{'\n'})
	}

	bs := bytes.SplitN(b, []byte{'\n'}, 3)
	line0, line1 := bs[0], bs[1]

	if prefix := "--- difftest.txt.orig"; !bytes.HasPrefix(line0, []byte(prefix)) {
		t.Errorf("diff: first line should start with `%s`\ngot: %s", prefix, line0)
	}

	if prefix := "+++ difftest.txt"; !bytes.HasPrefix(line1, []byte(prefix)) {
		t.Errorf("diff: second line should start with `%s`\ngot: %s", prefix, line1)
	}

	want := `@@ -1,2 +1,2 @@
 first
-second
+third
`

	if got := string(bs[2]); got != want {
		t.Errorf("diff: got:\n%s\nwant:\n%s", got, want)
	}
}

func TestReplaceTempFilename(t *testing.T) {
	diff := []byte(`--- /tmp/tmpfile1	2017-02-08 00:53:26.175105619 +0900
+++ /tmp/tmpfile2	2017-02-08 00:53:38.415151275 +0900
@@ -1,2 +1,2 @@
 first
-second
+third
`)
	want := []byte(`--- path/to/file.go.orig	2017-02-08 00:53:26.175105619 +0900
+++ path/to/file.go	2017-02-08 00:53:38.415151275 +0900
@@ -1,2 +1,2 @@
 first
-second
+third
`)
	// Check path in diff output is always slash regardless of the
	// os.PathSeparator (`/` or `\`).
	sep := string(os.PathSeparator)
	filename := strings.Join([]string{"path", "to", "file.go"}, sep)
	got, err := replaceTempFilename(diff, filename)
	if err != nil {
		t.Fatal(err)
	}
	if !bytes.Equal(got, want) {
		t.Errorf("os.PathSeparator='%s': replacedDiff:\ngot:\n%s\nwant:\n%s", sep, got, want)
	}
}
