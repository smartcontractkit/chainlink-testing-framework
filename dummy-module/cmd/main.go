package main

import (
	"fmt"
	dummy_module "github.com/smartcontractkit/chainlink-testing-framework/dummy-module"
)

func FuncTest() error {
	return nil
}

func main() {
	if err := dummy_module.FuncTest(); err != nil {
		panic(err)
	}
	fmt.Println("Hello!")
}
