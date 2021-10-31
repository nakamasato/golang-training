package blogposts

import (
	"errors"
	"io/fs"
	"reflect"
	"testing"
	"testing/fstest"
	// "github.com/nakamasato/blogposts"
)
func TestNewBlogPosts(t *testing.T) {
    const (
        firstBody = `Title: Post 1
Description: Description 1
Tags: tdd, go
---
Hello
World`
        secondBody = `Title: Post 2
Description: Description 2
Tags: rust, borrow-checker
---
B
L
M`
    )
    fs := fstest.MapFS{
        "hello world.md":  {Data: []byte(firstBody)},
        "hello-world2.md": {Data: []byte(secondBody)},
    }
    posts, err := NewPostsFromFS(fs)
    if err != nil {
        t.Fatal(err)
    }

    if len(posts) != len(fs) {
        t.Errorf("got %d posts, wanted %d posts", len(posts), len(fs))
    }

    assertPost(t, posts[0], Post{
        Title: "Post 1",
        Description: "Description 1",
        Tags:        []string{"tdd", "go"},
        Body: `Hello
World`,
    })
}

func assertPost(t *testing.T, got Post, want Post) {
    t.Helper()
    if !reflect.DeepEqual(got, want) {
        t.Errorf("got %+v, want %+v", got, want)
    }
}

type StubFailingFS struct {
}
func (s StubFailingFS) Open(name string) (fs.File, error) {
    return nil, errors.New("oh no, i always fail")
}
