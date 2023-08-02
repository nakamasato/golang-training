package main

import (
	"fmt"
	"math"
	"math/big"
	"strings"

	"github.com/google/uuid"
	"github.com/lithammer/shortuuid/v4"
)

func main() {
	for i := 1; i <= 10; i++ {
		shortuid := shortuuid.New()
		if len(shortuid) != 22 {
			fmt.Printf("shortuuid: %s, len: %d\n", shortuid, len(shortuid))
			break
		}
	}
	u := uuid.New()
	var num big.Int
	fmt.Println(strings.Replace(u.String(), "-", "", 4))
	num.SetString(strings.Replace(u.String(), "-", "", 4), 16)
	fmt.Println(num)
	alphabet := "23456789ABCDEFGHJKLMNPQRSTUVWXYZabcdefghijkmnopqrstuvwxyz"
	length := math.Ceil(math.Log(math.Pow(2, 128)) / math.Log(float64(len(alphabet))))
	fmt.Println(length)
}
