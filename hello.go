package main

// import "fmt"

const englishHelloPrefix = "Hello, "
const spanishHelloPrefix = "Hola, "

func Hello(name string, language string) string {
	if name == "" {
    	name = "World"
	}
	prefix := englishHelloPrefix
	if language == "Spanish" {
		prefix = spanishHelloPrefix
	}
	return prefix + name
}

// func main() {
//     fmt.Println(Hello("naka", "English"))
// }
