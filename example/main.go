package main

import (
	"fmt"

	"github.com/lixiangzhong/safemap"
)

func main() {
	m := safemap.New[string, string]()
	value := m.GetOrSet("a", "b")
	fmt.Println(value)
	// output: b
	value, ok := m.Get("a")
	fmt.Println(value, ok)
	// output: b true
}
