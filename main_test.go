package main

import (
	"bytes"
	"io/ioutil"
	"strings"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestParse_NoImports(t *testing.T) {
	_, rewritten := parse([]byte(""))
	require.False(t, rewritten)
}

func TestParse(t *testing.T) {
	type testcase struct {
		Description     string
		ActualFixture   string
		ExpectedFixture string
		NoGoFmt         bool
		GenCode         bool
	}
	testcases := []testcase{
		{
			Description:     "go-groups should not modify sorted import blocks",
			ActualFixture:   "valid_imports.txt",
			ExpectedFixture: "valid_imports.txt",
		},
		{
			Description:     "go-groups should merge multiple import groups",
			ActualFixture:   "extra_groups.txt",
			ExpectedFixture: "valid_imports.txt",
		},
		{
			Description:     "go-groups should group external imports",
			ActualFixture:   "external_groups_invalid.txt",
			ExpectedFixture: "external_groups.txt",
		},
		{
			Description:     "go-groups should sort external imports",
			ActualFixture:   "sort_external_groups_invalid.txt",
			ExpectedFixture: "sort_external_groups.txt",
		},
		{
			Description:     "go-groups should handle import prefixes",
			ActualFixture:   "prefixed_groups_invalid.txt",
			ExpectedFixture: "prefixed_groups.txt",
		},
		{
			Description:     "go-groups should group subdomain external imports",
			ActualFixture:   "subdomain_imports_invalid.txt",
			ExpectedFixture: "subdomain_imports.txt",
		},
		{
			Description:     "go-groups should preserve comments in imports",
			ActualFixture:   "commented_imports.txt",
			ExpectedFixture: "commented_imports.txt",
		},
		{
			Description:     "go-groups should group multi-line comments in imports",
			ActualFixture:   "multiline_comments_invalid.txt",
			ExpectedFixture: "multiline_comments.txt",
			NoGoFmt:         true,
		},
		{
			Description:     "go-groups should run gofmt over source code",
			ActualFixture:   "gofmt_invalid.txt",
			ExpectedFixture: "gofmt.txt",
		},
		{
			Description:     "go-groups should not run gofmt over source code",
			ActualFixture:   "gofmt_invalid.txt",
			ExpectedFixture: "gofmt_invalid.txt",
			NoGoFmt:         true,
		},
		{
			Description:     "go-groups should not run over generated code",
			ActualFixture:   "code_generated.txt",
			ExpectedFixture: "code_generated.txt",
		},
		{
			Description:     "go-groups should run over generated code",
			ActualFixture:   "code_generated.txt",
			ExpectedFixture: "code_generated_valid.txt",
			GenCode:         true,
		},
		{
			Description:     "go-groups order imports with lint comments",
			ActualFixture:   "import_with_lint_comment.txt",
			ExpectedFixture: "import_with_lint_comment.txt",
		},
	}
	var buf bytes.Buffer
	for _, testcase := range testcases {
		bytes := testdata(t, testcase.ActualFixture)
		expected := testdata(t, testcase.ExpectedFixture)
		err := processFile("", strings.NewReader(string(bytes)), &buf, !testcase.NoGoFmt, testcase.GenCode)

		if buf.String() != string(expected) {
			t.Logf("Input: \n%s\n\nOutput: \n%s\n\nExpected: \n%s\n\n", string(bytes), buf.String(), string(expected))
		}
		require.NoError(t, err)
		require.Equal(t, string(expected), buf.String(), testcase.Description)
		buf.Reset()
	}
}

func testdata(t *testing.T, str string) []byte {
	b, err := ioutil.ReadFile("testdata/" + str)
	require.NoError(t, err)
	return b
}
