# Meta

## [1. Why unit tests and how to make them work for you](https://quii.gitbook.io/learn-go-with-tests/meta/why)

[LondonGophers 12/12/2018: Chris James - How to not build legacy systems that everyone hates](https://www.youtube.com/watch?v=Kwtit8ZEK7U)

### Writing effective unit tests is a design problem

> it is desirable to have within your system self-contained, decoupled "units" centered around key concepts in your domain.

> I like to imagine these units as simple Lego bricks which have coherent APIs that I can combine with other bricks to make bigger systems. Underneath these APIs there could be dozens of things (types, functions et al) collaborating to make them work how they need to.

Unit tests are against "units" like I described. They were *never* about only being against a single class/function/whatever.

1. Refactoring
1. Unit tests
1. Unit design

### Why Test Driven Development (TDD)

TDD can help and force you to design well factored software iteratively, backed by tests to help future work as it arrives.

1. Write a small test for a small amount of desired behaviour
1. Check the test fails with a clear error (red)
1. Write the minimal amount of code to make the test pass (green)
1. Refactor
1. Repeat

## [2. Anti Pattern](https://quii.gitbook.io/learn-go-with-tests/meta/anti-patterns)

### Not doing TDD at all

One of the strengths of TDD is that it gives you **a formal process** to
- break down problems,
- understand what you're trying to achieve (red),
- get it done (green), then
- have a good think about how to make it right (blue/refactor).

### Misunderstanding the constraints of the refactoring step

Make a test pass in the refactoring step

### Having tests that won't fail (or, evergreen tests)

It's impossible in TDD, first step of which is

> Write a test, see it fail

### Useless assertions

Not a helpful message

> false was not equal to true

### Asserting on irrelevant detail

Rather than compare complex object, compare specific field!

### Lots of assertions within a single scenario for unit tests

They often creep in gradually, especially if test setup is complicated because you're reluctant to replicate the same horrible setup to assert on something else. Instead of this **you should fix the problems in your design which are making it difficult to assert on new things.**

- A helpful rule of thumb is to aim to make one assertion per test.
- take advantage of subtests

### Not listening to your tests

Remember
> TDD gives you the fastest feedback possible on your design

### Excessive setup, too many test doubles, etc.

**test double**
- **Dummy** objects are passed around but never actually used. Usually they are just used to fill parameter lists.
- **Fake** objects actually have working implementations, but usually take some shortcut which makes them not suitable for production (an `InMemoryTestDatabase` is a good example).
    <details><summary>Example</summary>

    [02-build-an-application/in_memory_player_store.go](../../02-build-an-application/in_memory_player_store.go)

    ```go
    type InMemoryPlayerStore struct {
        store map[string]int
    }

    func (i *InMemoryPlayerStore) RecordWin(name string) {
        i.store[name]++
    }

    func (i *InMemoryPlayerStore) GetPlayerScore(name string) int {
        return i.store[name]
    }

    func (i *InMemoryPlayerStore) GetLeague() League {
        var league []Player
        for name, wins := range i.store {
            league = append(league, Player{name, wins})
        }
        return league
    }
    ```
    </details>
- **Stubs** provide canned answers to calls made during the test, usually not responding at all to anything outside what's programmed in for the test.
    <details><summary>Example</summary>

    [02-build-an-application/testing.go](../../02-build-an-application/testing.go)

    ```go
    type StubPlayerStore struct {
        Scores   map[string]int
        WinCalls []string
        League   []Player
    }

    func (s *StubPlayerStore) GetPlayerScore(name string) int {
        score := s.Scores[name]
        return score
    }

    func (s *StubPlayerStore) RecordWin(name string) {
        s.WinCalls = append(s.WinCalls, name)
    }

    func (s *StubPlayerStore) GetLeague() League {
        return s.League
    }
    ```

    </details>
- **Spies** are **stubs** that also record some information based on how they were called. One form of this might be an email service that records how many messages it was sent.
    <details><summary>Example</summary>

    [02-build-an-application/CLI_test.go](../../02-build-an-application/CLI_test.go)

    ```go
    type GameSpy struct {
        StartCalled     bool
        StartCalledWith int
        BlindAlert      []byte

        FinishedCalled   bool
        FinishCalledWith string
    }

    func (g *GameSpy) Start(numberOfPlayers int, out io.Writer) {
        g.StartCalledWith = numberOfPlayers
        g.StartCalled = true
        _, err := out.Write(g.BlindAlert)
        if err != nil {
            log.Fatal(err)
        }
    }

    func (g *GameSpy) Finish(winner string) {
        g.FinishCalledWith = winner
    }
    ```

    [02-build-an-application/game_test.go](../../02-build-an-application/game_test.go)

    ```go
    type SpyBlindAlerter struct {
        alerts []scheduledAlert
    }

    func (s *SpyBlindAlerter) ScheduleAlertAt(duration time.Duration, amount int, to io.Writer) {
        s.alerts = append(s.alerts, scheduledAlert{duration, amount})
    }
    ```

    </details>
- **Mocks** are pre-programmed with expectations which form a specification of the calls they are expected to receive. They can throw an exception if they receive a call they don't expect and are checked during verification to ensure they got all the calls they were expecting.



Points:
- Leaky interfaces
- Think about the types of test doubles you use
    - **Mocks** are sometimes helpful, but they're extremely powerful and therefore easy to misuse. Try giving yourself the constraint of using **stubs** instead.
    - Verifying implementation detail with **spies** is sometimes helpful, but **try to avoid it**. Remember your implementation detail is usually not important, and you don't want your tests coupled to them if possible. Look to couple your tests to useful behaviour rather than incidental details.
    - [Start naming your test doubles correctly](https://quii.dev/Start_naming_your_test_doubles_correctly)
- [Consolidate dependencies](https://quii.gitbook.io/learn-go-with-tests/meta/anti-patterns#consolidate-dependencies)
    - many dependencies with HTTP handler example.

### Violating encapsulation
- the function being tested is only called from tests. Which is obviously a terrible outcome, and a waste of time.
- **In Go, consider your default position for writing tests as from the perspective of a consumer of your package.**
### Complicated table tests

### Summary

Most problems with unit tests can normally be traced to:
- Developers not following the TDD process
- Poor design

The good news is **TDD can help you improve your design skills** because as stated in the beginning:
**TDD's main purpose is to provide feedback on your design.**

## [3. Intro to generics](https://quii.gitbook.io/learn-go-with-tests/meta/intro-to-generics)

[generics](https://github.com/golang/go/issues/43651#issuecomment-776944155) will be included in version 1.18

you'll know how to write:
- A function that takes generic arguments
- A generic data-structure

-> https://gotipplay.golang.org/

**A function that takes a string or an integer? (or indeed, other things)**
- your argument as `interface{}` which means "anything".

```go
func AssertEqual(got, want interface{})

func AssertNotEqual(got, want interface{})
```

```go
func (is *I) Equal(a, b interface{})
```

-> no information at compile time as to what kinds of data we're dealing with.

**Our own test helpers with generics**


**Our own test helpers with generics**

[comparable](https://go.googlesource.com/proposal/+/refs/heads/master/design/43651-type-parameters.md#comparable-types-in-constraints)

```go
func AssertEqual[T comparable](got, want T) {
```

```go
func InterfaceyFoo(x, y interface{})
// and
func GenericFoo[T any](x, y T)
```

All valid for `InterfaceyFoo(apple, orange)`

Valid:
- GenericFoo(apple1, apple2)
- GenericFoo(orange1, orange2)
- GenericFoo(1, 2)
- GenericFoo("one", "two")
Not valid (fails compilation):
- GenericFoo(apple1, orange1)
- GenericFoo("1", 1)

**Next: Generic data types**

```go
type StackOfInts struct {
	values []int
}

func (s *StackOfInts) Push(value int) {
	s.values = append(s.values, value)
}
...
```

```go
type StackOfStrings struct {
	values []string
}

func (s *StackOfStrings) Push(value string) {
	s.values = append(s.values, value)
}
```

The code for both `StackOfStrings` and `StackOfInts` is almost identical.

pre-generic with `interface{}`:

```go
type Stack struct {
	values []interface{}
}

func (s *Stack) Push(value interface{}) {
	s.values = append(s.values, value)
}
```

Problems:
- I can now Push apples onto a stack of oranges.

**Generic data structures to the rescue**

```go
type Stack[T any] struct {
	values []T
}
```

https://gotipplay.golang.org/p/xAWcaMelgQV

```go
myStackOfStrings := new(Stack[string]) // specify the type
myStackOfInts := new(Stack[int]) // specify the type
```
