# Learn Go with tests

[![codecov](https://codecov.io/gh/nakamasato/go-practice/branch/main/graph/badge.svg?token=1RUXMSBB6N)](https://codecov.io/gh/nakamasato/go-practice)

https://quii.gitbook.io/

## Version

`go1.17.1`

## Commands

test:
- `go test -cover`
- `go test -v`
- `go test -race` race detector

## Install Go

- Modules

    ```bash
    go mod init <modulepath used "tmp">
    ```
- Debug: `Delve`
- Lint: [golangci-lint](https://golangci-lint.run/)

## Contents

- [Learn Go with tests](learn-go-with-tests)
    - [Go fundamentails](learn-go-with-tests/01-go-fundamentals)
    - [Build an application](learn-go-with-tests/02-build-an-application)
- [Pragmatic Cases](pragmatic-cases)
    - [Prometheus](pragmatic-cases/prometheus)
    - [kind cluster](pragmatic-cases/kind)
    - [Skaffold](pragmatic-cases/skaffold)
    - [Cobra](https://github.com/nakamasato/cobra-sample)
