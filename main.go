package main

import (
	"bufio"
	"bytes"
	"flag"
	"fmt"
	"os"
	"regexp"
	"sort"
	"strings"
)

const (
	exitBadFlags      = 1
	exitBadStdin      = 2
	exitInternalError = 3

	versionStr = "go-groups version 1.0.3 (2019-06-20)"
)

var (
	importStartRegex = regexp.MustCompile(`^\s*import\s*\(\s*$`)
	importEndRegex   = regexp.MustCompile(`^\s*\)\s*$`)

	// any whitespace + any unicode_letter_or_underscore + any unicode_letter_or_underscore_or_unicode number + any whitespace + quote + any + quote + any
	groupedImportRegex = regexp.MustCompile(`^\s*[\p{L}_]*[\s*[\p{L}_\p{N}]*\s*".*".*$`)
	externalImport     = regexp.MustCompile(`"([a-zA-Z0-9_]{1}[a-zA-Z0-9_-]{0,62}){1}(\.[a-zA-Z0-9_]{1}[a-zA-Z0-9_-]{0,62})*[\._]?\/([\p{L}_\-\p{N}]*)\/?.*"`)

	list    = flag.Bool("l", false, "list files whose formatting differs")
	write   = flag.Bool("w", false, "write result to (source) file instead of stdout")
	doDiff  = flag.Bool("d", false, "display diffs instead of rewriting files")
	version = flag.Bool("v", false, "display the version of go-groups")
)

func main() {
	flag.Usage = usage
	flag.Parse()

	if *version {
		fmt.Println(versionStr)
		os.Exit(0)
	}

	// stdin invocation
	if flag.NArg() == 0 {
		if *write {
			_, _ = fmt.Fprintln(os.Stderr, "error: cannot use -w with standard input")
			os.Exit(exitBadStdin)
		}
		if err := processFile("<standard input>", os.Stdin, os.Stdout); err != nil {
			_, _ = fmt.Fprintln(os.Stderr, "failed to parse stdin: "+err.Error())
			os.Exit(exitBadFlags)
		}
		return
	}

	for i := 0; i < flag.NArg(); i++ {
		path := flag.Arg(i)
		switch dir, err := os.Stat(path); {
		case err != nil:
			_, _ = fmt.Fprintln(os.Stderr, "no files matching '"+path+"': "+err.Error())
			os.Exit(exitBadFlags)
		case dir.IsDir():
			if err := walkDir(path); err != nil {
				_, _ = fmt.Fprintln(os.Stderr, "failed processing path "+path+": "+err.Error())
				os.Exit(exitInternalError)
			}
		default:
			_ = processFile(path, nil, os.Stdout)
		}
	}
}

func usage() {
	_, _ = fmt.Fprintf(os.Stderr, "usage: %s [flags] [path ...]\n", os.Args[0])
	flag.PrintDefaults()
}

type importGroup struct {
	lineStart int
	lineEnd   int

	lines []string
}

func parse(src []byte) (result []byte, rewritten bool) {
	groups := make([]importGroup, 0, 1)
	contents := make(map[int]string, 128)

	scanner := bufio.NewScanner(bytes.NewReader(src))

	insideImports := false
	var n int
	var group importGroup
	for n = 0; scanner.Scan(); n++ {
		line := scanner.Text()
		if insideImports {
			if importEndRegex.MatchString(line) {
				insideImports = false
				group.lineEnd = n
				groups = append(groups, group)
			} else if groupedImportRegex.MatchString(line) {
				group.lines = append(group.lines, line)
			}
		} else if importStartRegex.MatchString(line) {
			insideImports = true
			group = importGroup{
				lineStart: n,
			}
		} else {
			contents[n] = line
		}
	}

	// nothing to do
	if len(groups) == 0 {
		return []byte{}, false
	}

	for i, group := range groups {
		groups[i] = regroupImportGroups(group)
	}

	fileBytes := fixupFile(contents, n, groups)

	return fileBytes, true
}

func fixupFile(contents map[int]string, numLines int, groups []importGroup) []byte {
	buffer := bytes.NewBufferString("")
	for i := 0; i < numLines; i++ {
		line, ok := contents[i]
		if ok {
			buffer.WriteString(line)
			buffer.WriteString("\n")
			continue
		}
		for _, group := range groups {
			if group.lineStart == i {
				buffer.WriteString("import (\n")
				leadingWhitespace := true
				for _, line := range group.lines {
					if leadingWhitespace && strings.TrimSpace(line) == "" {
						// skip empty leading import lines
					} else {
						buffer.WriteString(line)
						buffer.WriteString("\n")
						leadingWhitespace = false
					}
				}
				buffer.WriteString(")\n")
				i = group.lineEnd
				break
			}
		}
	}
	return buffer.Bytes()
}

// regroupImportGroups iterates each line of the import group and sorts the imports
// standard library imports are grouped together and sorted alphabetically
// each second-level external import is grouped together (e.g github.com/pkg.* is one group)
// each of these second-level groups is discovered and sorted alphabetically
// then each import is matched with their group and the list of lines to be written is built up
func regroupImportGroups(group importGroup) importGroup {
	standardImports := make(Imports, 0, len(group.lines))

	sortedKeys := make([]string, 0)
	groupNames := make(map[string]Imports, 0)
	for _, importLine := range group.lines {
		matches := externalImport.FindStringSubmatch(importLine)

		if matches != nil && strings.ContainsAny(importLine, ".") {
			groupName := strings.Join(matches[1:], "")
			if groupNames[groupName] == nil {
				groupNames[groupName] = make(Imports, 0, 1)
				sortedKeys = append(sortedKeys, groupName)
			}
			groupNames[groupName] = append(groupNames[groupName], importLine)
		} else {
			standardImports = append(standardImports, importLine)
		}
	}
	sort.Sort(standardImports)
	sort.Strings(sortedKeys)

	group.lines = standardImports
	for _, groupName := range sortedKeys {
		imports := groupNames[groupName]
		sort.Sort(imports)

		group.lines = append(group.lines, "")
		group.lines = append(group.lines, imports...)
	}
	return group
}
