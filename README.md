# Golang Training

[![codecov](https://codecov.io/gh/nakamasato/golang-training/branch/main/graph/badge.svg?token=1RUXMSBB6N)](https://codecov.io/gh/nakamasato/golang-training)

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
    1. [Questions and answers](learn-go-with-tests/03-questions-and-answers)
    1. [Meta](learn-go-with-tests/04-meta)
1. [Pragmatic Cases](pragmatic-cases)
    1. [Prometheus](pragmatic-cases/prometheus)
    1. [kind cluster](pragmatic-cases/kind)
    1. [k8s client](pragmatic-cases/k8sclient)
    1. [Skaffold](pragmatic-cases/skaffold)
    1. [Cobra](https://github.com/nakamasato/cobra-sample)
    1. [Ginkgo](pragmatic-cases/ginkgo)
    1. [MySQL](pragamtic-cases/mysql)
    1. [String to Object](pragmatic-cases/string-to-object)

## References & readings
1. [Learn Go with Tests](https://quii.gitbook.io/)
1. [Advanced Testing with Go](https://speakerdeck.com/mitchellh/advanced-testing-with-go)
1. [よくわかるcontextの使い方](https://zenn.dev/hsaki/books/golang-context)
1. [Goでの並行処理を徹底解剖！](https://zenn.dev/hsaki/books/golang-concurrency)
