package main

import (
	"fmt"
	dummy_module "github.com/smartcontractkit/chainlink-testing-framework/dummy-module"
)

func main() {
	if err := dummy_module.FuncTest(0); err != nil {
		panic(err)
	}
	fmt.Println("Hello!")
}
