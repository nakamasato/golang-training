package blogposts

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"io/fs"
	"strings"
)

type Post struct {
    Title string
    Description string
    Tags []string
    Body string
}
func NewPostsFromFS(fileSystem fs.FS) ([]Post, error) {
    dir, err := fs.ReadDir(fileSystem, ".")
    if err != nil {
        return nil, err
    }
    var posts []Post
    for _, f := range dir {
        post, err := getPost(fileSystem, f.Name())
        if err != nil {
            return nil, err
        }
        posts = append(posts, post)
    }
    return posts, nil
}

func getPost(fileSystem fs.FS, filename string) (Post, error) {
    postFile, err := fileSystem.Open(filename)
    if err != nil {
        return Post{}, err
    }
    defer postFile.Close()
    return newPost(postFile)
}

const (
    titleSeparator       = "Title: "
    descriptionSeparator = "Description: "
    tagsSeparator = "Tags: "
    bodySeparator = "Body:"
)

func newPost(postFile io.Reader) (Post, error) {
    scanner := bufio.NewScanner(postFile)

    readMetaLine := func(tagName string) string {
        scanner.Scan()
        return strings.TrimPrefix(scanner.Text(), tagName)
    }

    title := readMetaLine(titleSeparator)
    description := readMetaLine(descriptionSeparator)
    tags := strings.Split(readMetaLine(tagsSeparator), ", ")

    return Post{
        Title: title,
        Description: description,
        Tags: tags,
        Body: readBody(scanner),
    }, nil
}

func readBody(scanner *bufio.Scanner) string {
    scanner.Scan() // ignore a line
    buf := bytes.Buffer{}
    for scanner.Scan() {
        fmt.Fprintln(&buf, scanner.Text())
    }
    return strings.TrimSuffix(buf.String(), "\n")
}
