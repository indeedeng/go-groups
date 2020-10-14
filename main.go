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

	versionStr = "go-groups version 1.1.3 (2020-10-14)"
)

var (
	importStartRegex = regexp.MustCompile(`^\s*import\s*\(\s*$`)
	importEndRegex   = regexp.MustCompile(`^\s*\)\s*$`)

	// any whitespace + any unicode_letter_or_underscore + any unicode_letter_or_underscore_or_unicode number + any whitespace + quote + any + quote + any.
	groupedImportRegex = regexp.MustCompile(`^\s*[\p{L}_\\.]*[\s*[\p{L}_\p{N}]*\s*".*".*$`)
	externalImport     = regexp.MustCompile(`"([a-zA-Z0-9_]{1}[a-zA-Z0-9_-]{0,62}){1}(\.[a-zA-Z0-9_]{1}[a-zA-Z0-9_-]{0,62})*[\._]?\/([\p{L}_\-\p{N}]*)\/?.*"`)
	// see https://golang.org/pkg/cmd/go/internal/generate/ for details.
	generatedRegex = regexp.MustCompile(`^// Code generated .* DO NOT EDIT.$`)

	list     = flag.Bool("l", false, "list files whose formatting differs")
	write    = flag.Bool("w", false, "write result to (source) file instead of stdout")
	doDiff   = flag.Bool("d", false, "display diffs instead of rewriting files")
	version  = flag.Bool("v", false, "display the version of go-groups")
	noFormat = flag.Bool("f", false, "disables the automatic gofmt style fixes")
	genCode  = flag.Bool("g", false, "include generated code in analysis")
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
		if err := processFile("<standard input>", os.Stdin, os.Stdout, !*noFormat, *genCode); err != nil {
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
			_ = processFile(path, nil, os.Stdout, !*noFormat, *genCode)
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

	lines []importLine
}

type importLine struct {
	line         string
	contentAbove string
	contentBelow string
}

func parse(src []byte) (result []byte, rewritten bool) {
	groups := make([]importGroup, 0, 1)
	contents := make(map[int]string, 128)

	scanner := bufio.NewScanner(bytes.NewReader(src))
	lines := make([]string, 0, 128)
	for scanner.Scan() {
		lines = append(lines, scanner.Text())
	}

	insideImports := false
	var n int
	var group importGroup
	scanner = bufio.NewScanner(bytes.NewReader(src))
	for n = 0; scanner.Scan(); n++ {
		line := scanner.Text()
		if insideImports {
			if importEndRegex.MatchString(line) {
				insideImports = false
				group.lineEnd = n
				groups = append(groups, group)
			} else if groupedImportRegex.MatchString(line) {
				importLine := importLine{
					line: line,
				}
				var above, below int
				for above = n - 1; above > 0; above-- {
					if groupedImportRegex.MatchString(lines[above]) || importStartRegex.MatchString(lines[above]) {
						above++
						break
					}
				}
				for below = n + 1; below < len(lines); below++ {
					if importEndRegex.MatchString(lines[below]) {
						below--
						break
					}
					// if we hit an import beneath us, assume it owns the non-import content above itself, not us
					if groupedImportRegex.MatchString(lines[below]) {
						below = n
						break
					}
				}
				if above != n {
					importLine.contentAbove = strings.Join(filterNewlines(lines[above:n]), "\n")
				}
				if below != n {
					importLine.contentBelow = strings.Join(filterNewlines(lines[n+1:below+1]), "\n")
				}
				group.lines = append(group.lines, importLine)
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

func filterNewlines(lines []string) []string {
	filtered := make([]string, 0, len(lines))
	for _, line := range lines {
		if strings.TrimSpace(line) != "" {
			filtered = append(filtered, line)
		}
	}
	return filtered
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
			if group.lineStart != i {
				continue
			}
			buffer.WriteString("import (\n")
			leadingWhitespace := true
			for _, importLine := range group.lines {
				if leadingWhitespace && strings.TrimSpace(importLine.line) == "" {
					// skip empty leading import lines
				} else {
					if importLine.contentAbove != "" {
						buffer.WriteString(importLine.contentAbove)
						buffer.WriteString("\n")
					}
					buffer.WriteString(importLine.line)
					buffer.WriteString("\n")
					if importLine.contentBelow != "" {
						buffer.WriteString(importLine.contentBelow)
						buffer.WriteString("\n")
					}
					leadingWhitespace = false
				}
			}
			buffer.WriteString(")\n")
			i = group.lineEnd
			break
		}
	}
	return buffer.Bytes()
}

// regroupImportGroups iterates each line of the import group and sorts the imports
// standard library imports are grouped together and sorted alphabetically
// each second-level external import is grouped together (e.g github.com/pkg.* is one group)
// each of these second-level groups is discovered and sorted alphabetically
// then each import is matched with their group and the list of lines to be written is built up.
func regroupImportGroups(group importGroup) importGroup {
	standardImports := make(Imports, 0, len(group.lines))

	sortedKeys := make([]string, 0)
	groupNames := make(map[string]Imports)
	for _, importLine := range group.lines {
		matches := externalImport.FindStringSubmatch(importLine.line)

		if matches != nil && strings.ContainsAny(importLine.line, ".") {
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

		group.lines = append(group.lines, importLine{})
		group.lines = append(group.lines, imports...)
	}
	return group
}
