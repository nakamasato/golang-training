# Learn Go with tests
https://quii.gitbook.io/

## Commands

test:
- `go test -cover`
- `go test -v`
- `go test -race` race detector

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
