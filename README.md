go-groups
=========

[![Go Report Card](https://goreportcard.com/badge/oss.indeed.com/go/go-groups)](https://goreportcard.com/report/oss.indeed.com/go/go-groups)
[![Build Status](https://travis-ci.com/indeedeng/go-groups.svg?branch=master)](https://travis-ci.org/indeedeng/go-groups)
[![GoDoc](https://godoc.org/oss.indeed.com/go/go-groups?status.svg)](https://godoc.org/oss.indeed.com/go/go-groups)
[![NetflixOSS Lifecycle](https://img.shields.io/osslifecycle/indeedeng/go-groups.svg)](OSSMETADATA)
[![GitHub](https://img.shields.io/github/license/indeedeng/go-groups.svg)](LICENSE)

# Project Overview

Command `go-groups` is a CLI tool to deterministically rewrite go import blocks to sort and re-group 
the imports.

`go-groups` is similar to `goimports`, but in addition to sorting imports 
lexigraphically, it also separates import blocks at a per-project level.
Like `goimports`, `go-groups` also runs gofmt and fixes any style/formatting
issues.

# Getting Started

The `go-groups` command can be installed by running:

```
$ go install oss.indeed.com/go/go-groups@latest
```

# Usage
```
$ go-groups -h
  usage: go-groups [flags] [path ...]
    -d    display diffs instead of rewriting files
    -f    disables the automatic gofmt style fixes
    -g    include generated code in analysis
    -l    list files whose formatting differs
    -v    display the version of go-groups
    -w    write result to (source) file instead of stdout
```

# Formatting Behavior

#### With this example source file input
```
import (
  "strings"
  "github.com/pkg/errors"
  "fmt"
  "github.com/gorilla/mux"
  "github.com/gorilla/csrf"
)
```

#### Running `go-groups` will produce
```
import (
  "fmt"
  "strings"
  
  "github.com/gorilla/csrf"
  "github.com/gorilla/mux"
  
  "github.com/pkg/errors"
)
```

#### Typical Workflow

Run `go-groups -w ./..` to rewrite and sort import groupings for go source files in a project.

# Asking Questions

For technical questions about `go-groups`, just file an issue in the GitHub tracker.

For questions about Open Source in Indeed Engineering, send us an email at
opensource@indeed.com

# Contributing

We welcome contributions! Feel free to help make `go-groups` better.

### Process

- Open an issue and describe the desired feature / bug fix before making
changes. It's useful to get a second pair of eyes before investing development
effort.
- Make the change. If adding a new feature, remember to provide tests that
demonstrate the new feature works, including any error paths. If contributing
a bug fix, add tests that demonstrate the erroneous behavior is fixed.
- Open a pull request. Automated CI tests will run. If the tests fail, please
make changes to fix the behavior, and repeat until the tests pass.
- Once everything looks good, one of the indeedeng members will review the
PR and provide feedback.

# Maintainers

The `oss.indeed.com/go/go-groups` module is maintained by Indeed Engineering.

While we are always busy helping people get jobs, we will try to respond to
GitHub issues, pull requests, and questions within a couple of business days.

# Code of Conduct

`oss.indeed.com/go/go-groups` is governed by the [Contributer Covenant v1.4.1](CODE_OF_CONDUCT.md)

For more information please contact opensource@indeed.com.

# License

The `oss.indeed.com/go/go-groups` module is open source under the [BSD-3-Clause](LICENSE) license.
