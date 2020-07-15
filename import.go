package main

import (
	"sort"
	"strings"
)

// Imports represents the list of imports in a given go file.
// This helper encapsulates the logic for sorting imports based on go-groups.
type Imports []importLine

var _ sort.Interface = (*Imports)(nil)

func (s Imports) Len() int {
	return len(s)
}

func (s Imports) Swap(i, j int) {
	s[i], s[j] = s[j], s[i]
}

func (s Imports) Less(i, j int) bool {
	s1 := strings.TrimSpace(s[i].line)
	if strings.ContainsAny(s1, " ") {
		s1 = strings.Join(strings.Split(s1, " ")[1:], " ")
	}
	s2 := strings.TrimSpace(s[j].line)
	if strings.ContainsAny(s2, " ") {
		s2 = strings.Join(strings.Split(s2, " ")[1:], " ")
	}
	return strings.Compare(s1, s2) < 0
}
