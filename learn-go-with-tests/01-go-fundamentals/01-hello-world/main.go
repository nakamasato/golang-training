package main

import "fmt"

const helloEnglish = "Hello"
const helloSpanish = "Ola"
const frenchHello = "Oi"
const defaultNameEnglish = "World"

func Hello(name string, language string) string {
	if len(name) == 0 {
		name = defaultNameEnglish
	}
	hello := GetHello(language)
	return hello + ", " + name + "!"
}

func GetHello(language string) string {
	switch language {
	case "French":
		return frenchHello
	case "Spanish":
		return helloSpanish
	default:
		return helloEnglish
	}
}

func main() {
	fmt.Println(Hello("John", "English"))
}
