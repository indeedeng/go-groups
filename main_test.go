package main

import (
	"io/ioutil"
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
	require.Equal(t, bytes, result)
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
	require.Equal(t, expected, result)
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
	require.Equal(t, expected, result)
}

func testdata(t *testing.T, str string) []byte {
	b, err := ioutil.ReadFile("testdata/" + str)
	require.NoError(t, err)
	return b
}
