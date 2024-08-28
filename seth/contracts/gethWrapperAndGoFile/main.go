package network_debug_sub_contract

import "math/big"

const constant = "package main"

var variable = "package main"

type SomeType struct {
	Int *big.Int
	Map map[string]string
}

type Interface interface {
	DoSomething()
}

func main() {
	_ = constant
	_ = variable
}
