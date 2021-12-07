# Ginkgo

https://onsi.github.io/ginkgo/

## Quickstart

1. Bootstrap a Ginkgo test suite.

    ```bash
    ginkgo bootstrap
    ```

    result: `book_suite_test.go` will be genarated.

1. Add specs to a suite.
    ```bash
    ginkgo generate book
    ```

    result: `book_test.go` will be genarated.

1. Write tests in `book_test.go`.
1. Run the tests.
    ```
    ginkgo
    ```

    ```
    Running Suite: Book Suite
    =========================
    Random Seed: 1638912576
    Will run 2 of 2 specs

    ••
    Ran 2 of 2 Specs in 0.000 seconds
    SUCCESS! -- 2 Passed | 0 Failed | 0 Pending | 0 Skipped
    PASS

    Ginkgo ran 1 suite in 1.378068683s
    Test Suite Passed
    ```

    or

    ```
    go test
    ```

    ```
    Running Suite: Book Suite
    =========================
    Random Seed: 1638912582
    Will run 2 of 2 specs

    ••
    Ran 2 of 2 Specs in 0.000 seconds
    SUCCESS! -- 2 Passed | 0 Failed | 0 Pending | 0 Skipped
    PASS
    ok      tmp/pragmatic-cases/ginkgo/book 0.137s
    ```

## Tips

### `Specify` Alias

```go
Describe("The foobar service", func() {
  Context("when calling Foo()", func() {
    Context("when no ID is provided", func() {
      Specify("an ErrNoID error is returned", func() {
      })
    })
  })
})
```

### `BeforeEach` and `AfterEach`

- Run before/after each spec.
- Share common state.

### Organizing Specs With Containers: `Describe` and `Context`

Example:

- `Describe`: `loading from JSON`
    - `Contest`:
        1. `when the JSON parses succesfully` -> normal tests
        1. `when the JSON fails to parse` with `BeforeEach` to make it fail to parse.

**Always initialize your variables in BeforeEach blocks** to avoid test pullution.

### `JustBeforeEach` and `JustAfterEach`

### Global Setup and Teardown: `BeforeSuite` and `AfterSuite`

### Documenting Complex `It`s: `By`

```go
By("Fetching a token and logging in")
```

### Pending Specs `P` or `X` before `Describe`, `Context`, `It`, and `Measure`
