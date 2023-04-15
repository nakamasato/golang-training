# Golang Training

[![codecov](https://codecov.io/gh/nakamasato/golang-training/branch/main/graph/badge.svg?token=1RUXMSBB6N)](https://codecov.io/gh/nakamasato/golang-training)

## Version

go1.20

## Basics

1. Install Go: https://go.dev/doc/install
    1. If you want to manage multiple versions, use [gvm](https://github.com/moovweb/gvm).
1. Run Go test:
    - `go test -cover`
    - `go test -v`
    - `go test -race` race detector
1. Init a module.
    ```bash
    go mod init <modulepath e.g. "tmp">
    ```
1. Tools:
    - Test:
        - BDD testing framework: [Ginkgo](https://onsi.github.io/ginkgo/)
        - Matcher/Assertion library: [Gomega](https://onsi.github.io/gomega/)
    - Debug: [Delve](https://github.com/go-delve/delve)
    - Lint: [golangci-lint](https://golangci-lint.run/)
    - VSCode extensions:
        - [vscode-ginkgo](https://marketplace.visualstudio.com/items?itemName=onsi.vscode-ginkgo)

## GitHub Actions

<details>

```yaml
      - uses: actions/checkout@v3

      - name: set up
        uses: actions/setup-go@v4
        with:
          go-version-file: go.mod
```

</details>

## Contents

1. [Learn Go with tests](learn-go-with-tests) (Official: https://quii.gitbook.io/)
    1. [Go fundamentails](learn-go-with-tests/01-go-fundamentals)
    1. [Build an application](learn-go-with-tests/02-build-an-application)
    1. [Questions and answers](learn-go-with-tests/03-questions-and-answers)
    1. [Meta](learn-go-with-tests/04-meta)
1. Pattern
    1. [Golang Functional Options Pattern](https://golang.cafe/blog/golang-functional-options-pattern.html)
1. [Pragmatic Cases](pragmatic-cases)
    1. Database
        1. [ent](pragmatic-cases/ent) (go1.17 is removed when upgrading to ent@v0.11.3 [#85](https://github.com/nakamasato/golang-training/pull/85))
        1. [MySQL](pragmatic-cases/mysql)
        1. [Migrate](pragmatic-cases/migrate)
            1. [Postgres](pragmatic-cases/migrate/postgres)
            1. [MySQL](pragmatic-cases/migrate/mysql)
        1. [atlas](pragmatic-cases/atlas)
    1. Kubernetes
        1. [kind cluster](pragmatic-cases/kind)
        1. [k8s client](pragmatic-cases/k8sclient) (needs go1.17 or later to use controller-runtime@v0.13.0 [#83](https://github.com/nakamasato/golang-training/pull/83))
        1. [Skaffold](pragmatic-cases/skaffold)
    1. Others
        1. [Prometheus](pragmatic-cases/prometheus)
        1. [Cobra](https://github.com/nakamasato/cobra-sample)
        1. [Ginkgo](pragmatic-cases/ginkgo)
        1. [String to Object](pragmatic-cases/string-to-object)
        1. [Opentelemetry](pragmatic-cases/opentelemetry)
        1. [gojsondiff](pragmatic-cases/gojsondiff)
## References & readings
1. [Learn Go with Tests](https://quii.gitbook.io/)
1. [Advanced Testing with Go](https://speakerdeck.com/mitchellh/advanced-testing-with-go)
1. [よくわかるcontextの使い方](https://zenn.dev/hsaki/books/golang-context)
1. [Goでの並行処理を徹底解剖！](https://zenn.dev/hsaki/books/golang-concurrency)
