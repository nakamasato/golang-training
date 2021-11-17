# Learn Go with tests
https://quii.gitbook.io/

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
    - [Expose Prometheus Metrics](pragmatic-cases/prometheus)
    - [Set up and tear down kind cluster](pragmatic-cases/kind)
    - [Deploy and delete resources by Skaffold](pragmatic-cases/skaffold)
