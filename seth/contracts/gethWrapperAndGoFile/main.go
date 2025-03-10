package network_debug_sub_contract

import "math/big"

// nolint
const constant = "package main"

// nolint
var variable = "package main"

type SomeType struct {
	Int *big.Int
	Map map[string]string
}

type Interface interface {
	DoSomething()
}

// nolint
func main() {
	_ = constant
	_ = variable
}
