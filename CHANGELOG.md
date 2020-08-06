# Changelog
All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.0.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [Unreleased]

## [1.1.1] - 2020-08-06
### Added
- Added -g flag to instruct go-groups to analyze and include generated code

### Changed
- go-groups now detects generated code and ignores generated code by default

## [1.1.0] - 2020-07-15
### Added
- Added golangci-lint linter, Makefile, and fixed linter warnings

### Changed
- go-groups now runs gofmt over the input files, unless disabled by the `-f` flag.

### Fixed
- Fixed go-groups stripping the dot from dot imports (thanks [@jdroot](https://github.com/jdroot))
- Fixed go-groups removing comments and other content from import blocks

## [1.0.3] - 2019-06-20
Initial public release
