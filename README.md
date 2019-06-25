go-groups
=========

# Project Overview

Command `go-groups` is a CLI tool to parse go import blocks, sort, and re-group 
the imports.

`go-groups` is similar to `goimports`, but in addition to sorting imports 
lexigraphically, it also separates import blocks at a per-project level.

# Getting Started

The `go-groups` command can be installed by running:

```
$ go get oss.indeed.com/go/go-groups
```

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

1. Run `go-groups -w ./..` to rewrite and sort import groupings
2. Run `go fmt -w ./..` to format your go source file

#### More info

`go-groups -h` to print usage help text and available flags

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
