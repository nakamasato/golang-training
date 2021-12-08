# Practice Go from Zero

[![codecov](https://codecov.io/gh/nakamasato/go-practice/branch/main/graph/badge.svg?token=1RUXMSBB6N)](https://codecov.io/gh/nakamasato/go-practice)

## Version

`go1.17.1`

## Basics

1. Install Go: https://go.dev/doc/install
1. Run Go test:
    - `go test -cover`
    - `go test -v`
    - `go test -race` race detector
1. Init a module.
    ```bash
    go mod init <modulepath used "tmp">
    ```
1. Tools:
    - Test:
        - BDD testing framework: [Ginkgo](https://onsi.github.io/ginkgo/)
        - Matcher/Assertion library: [Gomega](https://onsi.github.io/gomega/)
    - Debug: [Delve](https://github.com/go-delve/delve)
    - Lint: [golangci-lint](https://golangci-lint.run/)
    - VSCode extensions:
        - [vscode-ginkgo](https://marketplace.visualstudio.com/items?itemName=onsi.vscode-ginkgo)

## Contents

1. [Learn Go with tests](learn-go-with-tests) (Official: https://quii.gitbook.io/)
    1. [Go fundamentails](learn-go-with-tests/01-go-fundamentals)
    1. [Build an application](learn-go-with-tests/02-build-an-application)
1. [Pragmatic Cases](pragmatic-cases)
    1. [Prometheus](pragmatic-cases/prometheus)
    1. [kind cluster](pragmatic-cases/kind)
    1. [Skaffold](pragmatic-cases/skaffold)
    1. [Cobra](https://github.com/nakamasato/cobra-sample)
    1. [Ginkgo](pragmatic-cases/ginkgo)
