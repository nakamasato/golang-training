# Build an application

- GET `/players/<player>`
- POST `/players/<player>`

## Run the app

```
go build
```

```
curl http://localhost:5000/players/Pepper
```

```
curl -X POST http://localhost:5000/players/Pepper
```

## Server


### Step 1: Implement some function and make handler with HandlerFunc

![](docs/step1.drawio.svg)

```go
func ListenAndServe(addr string, handler Handler) error

type Handler interface {
    ServeHTTP(ResponseWriter, *Request)
}
```

[HandlerFunc](https://pkg.go.dev/net/http#HandlerFunc)

The `HandlerFunc` type is an adapter to allow the use of ordinary functions as HTTP handlers.

```go
func PlayerServer(w http.ResponseWriter, r *http.Request) {
    player := strings.TrimPrefix(r.URL.Path, "/players/")
    fmt.Fprint(w, GetPlayerScore(player))
}

func GetPlayerScore(name string) string {
    if name == "Pepper" {
        return "20"
    }
    if name == "Floyd" {
        return "10"
    }
    return ""
}
```

### Step 2: Implement Handler

![](docs/step2.drawio.svg)

```go
type PlayerServer struct {
    store PlayerStore
}

func (p *PlayerServer) ServeHTTP(w http.ResponseWriter, r *http.Request) {
    player := strings.TrimPrefix(r.URL.Path, "/players/")
    fmt.Fprint(w, p.store.GetPlayerScore(player))
}
```

## Step 3: Routing with [ServeMux](https://golang.org/pkg/net/http/#ServeMux)

```go
router := http.NewServeMux()
router.Handle("/league", http.HandlerFunc(funcA))
router.Handle("/players/", http.HandlerFunc(funcB))

func funcA(w http.ResponseWriter, r *http.Request) {
    // logic for league
}

func funcB(w http.ResponseWriter, r *http.Request) {
    // logic for player
}
```

## Step 4: Set Routing in NewPlayerServer

```go
func NewPlayerServer(store PlayerStore) *PlayerServer {
    p := &PlayerServer{
        store,
        http.NewServeMux(),
    }
    p.router.Handle("/league", http.HandlerFunc(p.leagueHandler))
    p.router.Handle("/players/", http.HandlerFunc(p.playersHandler))
    return p
}
```

## Step 5: Replace `ServeHTTP()` with `http.Handler` in `PlayerServer` by ***embedding***

https://pkg.go.dev/net/http#HandlerFunc.ServeHTTP

![](docs/step5.drawio.svg)

```go
type PlayerServer struct {
    store PlayerStore
	http.Handler // embedding: our PlayerServer now has all the methods that http.Handler has, which is just ServeHTTP.
}
```

[***embedding***](https://golang.org/doc/effective_go#embedding): *Go does not provide the typical, type-driven notion of subclassing, but it does have the ability to “borrow” pieces of an implementation by embedding types within a struct or interface.*


```go
func NewPlayerServer(store) *PlayerServer {
    p := new(PlayerServer)
    router := http.NewServerMux()
    ..
    p.Hander = router
    return p
}
```

We can call `ServeHTTP`

```go
store := StubPlayerStore{}
server := NewPlayerServer(&store)
server.ServeHTTP(w, r)
```

**Embedding** is a very interesting language feature. You can use it with interfaces to compose new interfaces. But you need to be careful with misuse (unintended exposion of methods of the embedded class).

```go
type Animal interface {
    Eater
    Sleeper
}
```

**SideNote**: about http.Handler & ServeMux

The process of the following code is nearly the same:
1. set the handler with `Handle` func
    ```go
    http.Handle("/any/", anyHandler)
    http.ListenAndServe(":8080", nil)
    ```
1. set the handler with `ServeMux`
    ```go
    mux := http.NewServeMux()
    mux.Handle("/any/", anyHandler)
    http.ListenAndServe(":8080", mux)
    ```

## Step 6: JSON encode and decode

- Use `encoding/json` package
- Encode `NewEncoder` with `io.Writer` (`http.ResponseWriter` in the example)
    ```go
    json.NewEncoder(w).Encode(leagueTable)
    ```
- Decode `newDecoder` with `io.Reader` (`response.Body` <- `*bytes.Buffer`)
    ```go
    obj := []Object
    json.NewDecoder(r).Decode(&obj)
    ```

## Step 7: Persist data with JSON file

- `FileSystemPlayerStore` with `database` as `io.Reader`
- Use `strings.NewReader`to create database with JSON string.
    - `newDecoder` to read json as object.

-> this implementation cannot read same thing twice once the Reader reachs the end

## Step 8: ReadSeeker: Seek(offset, whence) enables us to read multiple times

[ReadSeeker](https://golang.org/pkg/io/#ReadSeeker) interface

```go
type ReadSeeker struct {
    Reader
    Seeker
}
```

You can set offset and whence with:

```go
f.database.Seek(0, 0)
```

Luckily, `string.NewReader` also implements `ReadSeeker`.

## Step 9: ReadWriteSeeker

`strings.Reader` does not implement `ReadWriteSeeker`

- Create a temporary file: `*os.File` implements `ReadWriteSeeker`
- Use a [filebuffer](https://github.com/mattetti/filebuffer) library [Mattetti](https://github.com/mattetti): implements the interface

Replace `strings.Reader` with `*os.File` returned by `ioutil.TempFile("", "db")`.

```go
json.NewEncoder(f.database).Encode(league)
```
## [Step 10: More refactoring and performance concerns](https://quii.gitbook.io/learn-go-with-tests/build-an-application/io#more-refactoring-and-performance-concerns)

**Problem**: Every time someone calls `GetLeague()` or `GetPlayerScore()` we are reading the **entire** file and parsing it into JSON. -> Wasteful
**Solution**: Read the whole file only when starting up

1. Add `league` member variable to `FileSystemPlayerStore`.
    ```diff
     type FileSystemPlayerStore struct {
            database io.ReadWriteSeeker
    +       league   League
     }
    ```
1. `GetLeague` just returns from the member variable

    ```go
    func (f *FileSystemPlayerStore) GetLeague() League {
        return f.league
    }
    ```
1. `RecordWins` updates the member variable `league` and write with `f.database` (`io.ReadWriteSeeker`)
    ```diff
    -func (f FileSystemPlayerStore) RecordWin(name string) {
    -       league := f.GetLeague()
    -       player := league.Find(name)
    +func (f *FileSystemPlayerStore) RecordWin(name string) {
    +       player := f.league.Find(name)
    +
            if player != nil {
                    player.Wins++
            } else {
    -               league = append(league, Player{Name: name, Wins: 1})
    +               f.league = append(f.league, Player{Name: name, Wins: 1})
            }
    -       f.database.Seek(0, 0)
    -       json.NewEncoder(f.database).Encode(league)
    +       json.NewEncoder(f.database).Encode(f.league)
     }
    ```
1. Use `NewFileSystemPlayerStore(database)` to initialize.
    ```diff
    +func NewFileSystemPlayerStore(database io.ReadWriteSeeker) *FileSystemPlayerStore {
    +       database.Seek(0, 0)
    +       league, _ := NewLeague(database)
    +       return &FileSystemPlayerStore{
    +               database: database,
    +               league:   league,
    +       }
    +}
    ```

## [Step 11: Separate out the concern of the kind of data we write, from the writing](https://quii.gitbook.io/learn-go-with-tests/build-an-application/io#another-problem)

Problem: If new data is smaller than old data, the old data would remain at the end of the new data. e.g. Write `12345` -> Write `abc` -> Read `abc45` <- Wrong!!) -> Separate out the concern of the kind of **data we write**, from **the writing**

1. Introduce `tape.go`: encapsulate our **"when we write we go from the beginning" functionality**.
    ```go
    package main

    import "io"

    type tape struct {
           file io.ReadWriteSeeker
    }

    func (t *tape) Write(p []byte) (n int, err error) {
           t.file.Seek(0, 0)
           return t.file.Write(p)
    }
    ```
1. Update database of FileSystemPlayerStore from `io.ReadWriteSeeker` to `io.Writer`.
    ```diff
     type FileSystemPlayerStore struct {
    -       database io.ReadWriteSeeker
    +       database io.Writer
            league   League
     }
    ```
1. Update constructor of `FileSystemPlayerStore` to use `type`.
    ```diff
    @@ -14,7 +14,7 @@ func NewFileSystemPlayerStore(database io.ReadWriteSeeker) *FileSystemPlayerStor
            database.Seek(0, 0)
            league, _ := NewLeague(database)
            return &FileSystemPlayerStore{
    -               database: database,
    +               database: &tape{database},
                    league:   league,
            }
     }
    ```

Separation of concern completed!

## [Step 12 Enable to truncate the old data](https://quii.gitbook.io/learn-go-with-tests/build-an-application/io#write-the-test-first-4)

1. Add test case: Write `12345` -> Write `abc` -> Read `abc` (`file` -> `&tape{file}` -> `tape.Write([]byte("abc"))`) `tape_test.go`
    ```go
    package main

    import (
           "io/ioutil"
           "testing"
    )

    func TestTape_Write(t *testing.T) {
           file, clean := createTempFile(t, "12345")
           defer clean()

           tape := &tape{file}

           tape.Write([]byte("abc"))

           file.Seek(0, 0)
           newFileContents, _ := ioutil.ReadAll(file)

           got := string(newFileContents)
           want := "abc"

           if got != want {
                   t.Errorf("got %q want %q", got, want)
           }
    }
    ```
1. `os.File` has a truncate function. Use the type instead of `io.ReadWriteSeeker`
    ```diff
     package main

    -import "io"
    +import (
    +       "os"
    +)

     type tape struct {
    -       file io.ReadWriteSeeker
    +       file *os.File
     }

     func (t *tape) Write(p []byte) (n int, err error) {
    +       t.file.Truncate(0)
            t.file.Seek(0, 0)
            return t.file.Write(p)
     }
    ```

    For more details about `io.Reader`, `io.Writer`, and `*os.File`, check [Appendix](##Appendix)

1. Fix the compile errors (Change **data's type** from `io.ReadWriteSeeker` to `*os.File`).

    `file_system_store_test.go`:

    ```diff
    +++ b/learn-go-with-tests/02-build-an-application/file_system_store_test.go
    @@ -1,13 +1,12 @@
     package main

     import (
    -       "io"
            "io/ioutil"
            "os"
            "testing"
     )

    -func createTempFile(t testing.TB, initialData string) (io.ReadWriteSeeker, func()) {
    +func createTempFile(t testing.TB, initialData string) (*os.File, func()) {
            t.Helper()
            tmpfile, err := ioutil.TempFile("", "db")
            if err != nil {
    ```

    `file_system_store.go`:

    ```diff
    +++ b/learn-go-with-tests/02-build-an-application/file_system_store.go
    @@ -3,6 +3,7 @@ package main
     import (
            "encoding/json"
            "io"
    +       "os"
     )

     type FileSystemPlayerStore struct {
    @@ -10,7 +11,7 @@ type FileSystemPlayerStore struct {
            league   League
     }

    -func NewFileSystemPlayerStore(database io.ReadWriteSeeker) *FileSystemPlayerStore {
    +func NewFileSystemPlayerStore(database *os.File) *FileSystemPlayerStore {
            database.Seek(0, 0)
            league, _ := NewLeague(database)
            return &FileSystemPlayerStore{
    ```

## [Step 13 Small refactor](https://quii.gitbook.io/learn-go-with-tests/build-an-application/io#one-other-small-refactor)

Move `json.NewEncoder(f.database)` from `WriteWin` to `NewFileSystemPlayerStore`. No need to encode every time we write.

`file_system_store.go`:

```diff
 import (
        "encoding/json"
-       "io"
        "os"
 )

 type FileSystemPlayerStore struct {
-       database io.Writer
+       database *json.Encoder
        league   League
 }
@@ -15,7 +14,7 @@ func NewFileSystemPlayerStore(database *os.File) *FileSystemPlayerStore {
        database.Seek(0, 0)
        league, _ := NewLeague(database)
        return &FileSystemPlayerStore{
-               database: &tape{database},
+               database: json.NewEncoder(&tape{database}),
                league:   league,
        }
 }
@@ -40,5 +39,5 @@ func (f *FileSystemPlayerStore) RecordWin(name string) {
        } else {
                f.league = append(f.league, Player{Name: name, Wins: 1})
        }
-       json.NewEncoder(f.database).Encode(f.league)
+       f.database.Encode(f.league)
 }
```

`main.go`:

```diff
 package main

 import (
+       "encoding/json"
        "log"
        "net/http"
        "os"
@@ -14,7 +15,7 @@ func main() {
                log.Fatalf("problem opening %s %v", dbFileName, err)
        }

-       store := &FileSystemPlayerStore{db, League{}}
+       store := &FileSystemPlayerStore{json.NewEncoder(db), League{}}
        server := NewPlayerServer(store)
        // server := NewPlayerServer(NewInMemoryPlayerStore())
```


## Reference

- https://www.yunabe.jp/docs/golang_io.html

## Appendix

1. [io.Reader](https://pkg.go.dev/io#Reader): interface with `Read(p []byte) (n int, err error)` method. Usually not directly used.
1. [io.Writer](https://pkg.go.dev/io#Writer): interface with `Write(p []byte) (n int, err error)` method. Ususally not directly used.
1. [io.ReadSeeker](https://pkg.go.dev/io#ReadSeeker): interface with `Reader` and `Seeker`.
1. [io.ReadWriteSeeker](https://pkg.go.dev/io#ReadWriteSeeker): interface with `Reader`, `Writer`, and `Seeker`.
1. `os.Open`: Open a file and return `*os.File`, which can be used for `io.Reader` and `io.Writer`.
