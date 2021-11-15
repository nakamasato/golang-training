package main

import (
	"log"
	"os"
	"tmp/learn-go-with-tests/17-reading-files/blogposts"
)


func main() {
    posts, err := blogposts.NewPostsFromFS(os.DirFS("posts"))
    if err != nil {
        log.Fatal(err)
    }
    log.Println(posts)
}
