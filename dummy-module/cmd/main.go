package main

import "fmt"

func FuncTest(a int) error {
	return nil
}

func main() {
	if err := FuncTest(0); err != nil {
		panic(err)
	}
	fmt.Println("Hello!")
}
