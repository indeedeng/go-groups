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

func TestParse_ValidGroupingSortImports(t *testing.T) {
	bytes := testdata(t, "valid_imports.txt")
	result, rewritten := parse(bytes)

	require.True(t, rewritten)
	require.Equal(t, string(bytes), string(result))
}

func TestParse_ExtraGroups(t *testing.T) {
	bytes := testdata(t, "extra_groups.txt")
	result, rewritten := parse(bytes)

	require.True(t, rewritten)
	require.Equal(t, testdata(t, "valid_imports.txt"), result)
}

func TestParse_ExternalGroups(t *testing.T) {
	bytes := testdata(t, "external_groups_invalid.txt")
	expected := testdata(t, "external_groups.txt")
	result, rewritten := parse(bytes)

	require.True(t, rewritten)
	require.Equal(t, string(expected), string(result))
}

func TestParse_SortingExternalGroups(t *testing.T) {
	bytes := testdata(t, "sort_external_groups_invalid.txt")
	expected := testdata(t, "sort_external_groups.txt")
	result, rewritten := parse(bytes)

	require.True(t, rewritten)
	require.Equal(t, expected, result)
}

func TestParse_PrefixedGroups(t *testing.T) {
	bytes := testdata(t, "prefixed_groups_invalid.txt")
	expected := testdata(t, "prefixed_groups.txt")
	result, rewritten := parse(bytes)

	require.True(t, rewritten)
	require.Equal(t, expected, result)
}

func TestParse_SubdomainImports(t *testing.T) {
	bytes := testdata(t, "subdomain_imports_invalid.txt")
	expected := testdata(t, "subdomain_imports.txt")
	result, rewritten := parse(bytes)

	require.True(t, rewritten)
	require.Equal(t, string(expected), string(result))
}

func TestParse_CommentedImports(t *testing.T) {
	bytes := testdata(t, "commented_imports.txt")
	result, rewritten := parse(bytes)
	require.True(t, rewritten)
	require.Equal(t, string(bytes), string(result))
}

func TestParse_Multiline_CommentedImports(t *testing.T) {
	bytes := testdata(t, "multiline_comments_invalid.txt")
	expected := testdata(t, "multiline_comments.txt")
	result, rewritten := parse(bytes)
	require.True(t, rewritten)
	require.Equal(t, string(expected), string(result))
}

func TestParse_gofmt(t *testing.T) {
	b := testdata(t, "gofmt_invalid.txt")
	expected := testdata(t, "gofmt.txt")

	var buf bytes.Buffer
	err := processFile("", strings.NewReader(string(b)), &buf, true, false)
	require.NoError(t, err)

	require.Equal(t, string(expected), buf.String())
}

func TestParse_gofmt_disabled(t *testing.T) {
	b := testdata(t, "gofmt_invalid.txt")
	expected := testdata(t, "gofmt_invalid.txt")

	var buf bytes.Buffer
	err := processFile("", strings.NewReader(string(b)), &buf, false, false)
	require.NoError(t, err)

	require.Equal(t, string(expected), buf.String())
}

func TestParse_generated_code(t *testing.T) {
	b := testdata(t, "code_generated.txt")
	expected := testdata(t, "code_generated.txt")

	var buf bytes.Buffer
	err := processFile("", strings.NewReader(string(b)), &buf, true, false)
	require.NoError(t, err)

	require.Equal(t, string(expected), buf.String())
}

func TestParse_generated_code_enabled(t *testing.T) {
	b := testdata(t, "code_generated.txt")
	expected := testdata(t, "code_generated_valid.txt")

	var buf bytes.Buffer
	err := processFile("", strings.NewReader(string(b)), &buf, true, true)
	require.NoError(t, err)

	require.Equal(t, string(expected), buf.String())
}

func testdata(t *testing.T, str string) []byte {
	b, err := ioutil.ReadFile("testdata/" + str)
	require.NoError(t, err)
	return b
}
