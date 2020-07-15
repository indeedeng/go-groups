# Changelog
All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.0.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [Unreleased]

## [1.1.0] - 07-15-20
### Added
- Added golangci-lint linter, Makefile, and fixed linter warnings

### Changed
- go-groups now runs gofmt over the input files, unless disabled by the `-f` flag.

### Fixed
- Fixed go-groups stripping the dot from dot imports (thanks [@jdroot](https://github.com/jdroot))
- Fixed go-groups removing comments and other content from import blocks

## [1.0.3] - 2019-06-20
Initial public release
