package main

import "fmt"

var features = []string{
	"Free Feature #1",
	"Free Feature #2",
}

func main() {
	for _, f := range features {
		fmt.Println(">", f)
	}
}
