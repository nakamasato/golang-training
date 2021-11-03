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
## [Hello World](https://quii.gitbook.io/learn-go-with-tests/go-fundamentals/hello-world)

test

```go
func TestXXX(t *testing.T) {
    if got != want {
        t.Errorf("xxx")
    }
}
```

## [Integers](https://quii.gitbook.io/learn-go-with-tests/go-fundamentals/integers)

```go
Add(a, b int)
```

## [Iteration](https://quii.gitbook.io/learn-go-with-tests/go-fundamentals/iteration)

- [Benchmarks](https://pkg.go.dev/testing#hdr-Benchmarks)
    ```go
    func BenchmarkXxx(*testing.B)
    ```

    Example:

    ```go
    func BenchmarkRandInt(b *testing.B) {
        for i := 0; i < b.N; i++ {
            rand.Int()
        }
    }
    ```
- [Examples](https://go.dev/blog/examples) + [Examples](https://pkg.go.dev/testing#hdr-Examples)

    ```go
    func ExampleRepeat() {
        repeated := Repeat("ab", 3)
        fmt.Println(repeated)
        // Output: ababab
    }
    ```

## [Arrays and slices](https://quii.gitbook.io/learn-go-with-tests/go-fundamentals/arrays-and-slices)

- `range`
    ```go
    for _, s := range Numbers {
        fmt.Println(s)
    }
    ```
- `append`
    ```go
    SummedNumbers = append(SummedNumbers, Sum(numbers))
    ```
- `make`
    ```go
    sums := make([]int, lengthOfNumbers)
    ```
- [reflect.DeepEqual](https://golang.org/pkg/reflect/#DeepEqual)

    ```go
    want := []int{1, 2}
    got := []int{2, 3}
    if !reflect.DeepEqual(got, want) {
        ...
    }
    ```

## [Structs, methods & interfaces](https://quii.gitbook.io/learn-go-with-tests/go-fundamentals/structs-methods-and-interfaces)

- Struct
    ```go
    type Circle struct {
        Radius float64
    }
- Method
    ```go
    func (c Circle) Area() {
        return c.Radius * c.Radius * math.Pi
    }
- Interface
    ```go
    type Shape interface {
        Area()
    }
    ```
    - In Go interface resolution is implicit. All the struct that has `Area()` is recognized as `Shape`.
    - Using interfaces to declare **only what you need** is very important in software design

- Table driven test
    ```go
    areaTests := []struct {
        name    string
        shape   Shape
        hasArea float64
    }[
        {name: "test", shape: Triangle(Base: 40.0, Height: 10.0), hasArea: 200.0},
        ...
    ]
    for _, tt range areaTests {
        t.Run(tt.name, func(t *testing.T) {...})
    }
    ```

## [Pointers & errors](https://quii.gitbook.io/learn-go-with-tests/go-fundamentals/pointers-and-errors)

- In Go, **when you call a function or a method the arguments are *copied*.**
    ❌:
    ```go
    func (w Wallet) Deposit(amount int)  { // as w is a copy of whatever we called the method from
    	w.balance += amount
    }
    ```
    ⭕:
    ```go
    func (w *Wallet) Deposit(amount int)  {
    	w.balance += amount
    }
    ```
- the code above using (*w) is absolutely valid. However, the makers of Go deemed this notation cumbersome, so the language permits us to write w.balance, without an explicit dereference. ([automatic dereference](https://golang.org/ref/spec#Method_values))
    ```go
    func (w *Wallet) Balance() int {
    return (*w).balance // we can write w.balance!
    }
    ```

- `type Bitcoin int`: You can add domain specific functionality on top of existing types!!

    ```go
    func (b Bitcoin) String() string {
        return fmt.Sprintf("%d BTC", b)
    }
    ```

    -> We can use `"got %s want %s, got, want"`

- `t.Fatal`: will stop the test if it is called.
- `errcheck`: https://github.com/kisielk/errcheck

    ```
    go get -u github.com/kisielk/errcheck
    ```

    ```
    errcheck .
    wallet_test.go:15:24:   wallet.Withdraw(Bitcoin(10))
    ```

    If you can't find the command: https://githubmemory.com/repo/kisielk/errcheck/issues/194

## [](https://quii.gitbook.io/learn-go-with-tests/go-fundamentals/structs-methods-and-interfaces)

## [Dependency Injection](https://quii.gitbook.io/learn-go-with-tests/go-fundamentals/dependency-injection)

## [Mocking](https://quii.gitbook.io/learn-go-with-tests/go-fundamentals/mocking)

- use interface
- [**test double**](https://martinfowler.com/bliki/TestDouble.html): *Test Double is a generic term for any case where you replace a production object for testing purposes.*

## [Concurrency](https://quii.gitbook.io/learn-go-with-tests/go-fundamentals/concurrency)

**Goroutine**: run in a separate process

```go
go func(u string) {
    ...
}(url)
```

**Channel**: help organize and control the communication between thedifferent processes, allowing us to avoid a *race condition* bug.

```go
resultChannel <- result{u, wc(u)}
...
r := <-resultChannel
```

## [Select](https://quii.gitbook.io/learn-go-with-tests/go-fundamentals/select)

In the mocking and dependency injection chapters, we covered how ideally we don't want to be relying on external services to test our code because they can be
- Slow
- Flaky
- Can't test edge cases

**defer**: By prefixing a function call with defer it will now call that function at the end of the containing function.


```go
func Racer(a, b string) (winner string) {
    select { // get first one by running simultaneously
    case <-ping(a):
        return a
    case <-ping(b):
        return b
    }
}
```

You can add the following code to set timeout:

```go
case <-time.After(10 * time.Second):
```

## [Reflection](https://quii.gitbook.io/learn-go-with-tests/go-fundamentals/reflection)

**Reflection** in computing is the ability of a program to examine its own structure, particularly through types; it's a form of metaprogramming. It's also a great source of confusion.

- This can be quite clumsy and difficult to read and is generally less performant (as you have to do checks at runtime).
- In short **only use reflection** if you **really need to**.

- [reflect package](https://pkg.go.dev/reflect)
- [The Laws of Reflection](https://go.dev/blog/laws-of-reflection)

**Interface**: You may come across scenarios though where you want to write a function where you don't know the type at compile time. -> Go lets us get around this with the type interface{} which you can think of as just any type.

For more details: [reflection](reflection)

## [Sync](https://quii.gitbook.io/learn-go-with-tests/go-fundamentals/sync)

To make a function synchronized, we can add a lock.

> A Mutex is a mutual exclusion lock. The zero value for a Mutex is an unlocked mutex.

- **bad**

```go
import sync

type Counter struct {
	mu sync.Mutex
	num int
}
func (c *Counter) Inc()  {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.num++
}
```

This looks nice but while programming is a hugely subjective discipline, this is bad and wrong.

[sync.Mutex](https://pkg.go.dev/sync#Mutex)

> A Mutex must not be copied after first use.

-> pass in a pointer

**channels and goroutines** vs **mutex**

- Use channels when passing ownership of data
- Use mutexes for managing state
