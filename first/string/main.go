package main

import "fmt"

func main() {
	name := "naka"
	firstChar := name[0]
	secondChar := name[1]

	fmt.Println(firstChar) // 110
	fmt.Println(secondChar) // 97
	fmt.Println(string(firstChar)) // n
	fmt.Println(string(secondChar)) // a

	test := "abcdefg"

    // string 型 → []byte 型
    b := []byte(test)
    fmt.Print(b) // [97 98 99 100 101 102 103]

    // []byte 型 → string 型
    s := string(b)
    fmt.Print(s) // abcdefg
}
