package main

import "fmt"

func FuncTest() error {
	return nil
}

func main() {
	if err := FuncTest(); err != nil {
		panic(err)
	}
	fmt.Println("Hello!")
}
