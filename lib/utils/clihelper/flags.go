package clihelper

import "fmt"

// BoolFlag is a custom flag type to capture if a flag has been explicitly set.
type BoolFlag struct {
	IsSet bool
	Value bool
}

// Set is the method to set the flag value, part of the flag.Value interface.
func (f *BoolFlag) Set(s string) error {
	f.Value = s == "true"
	f.IsSet = true
	return nil
}

// String is the method to format the flag's value, part of the flag.Value interface.
func (f *BoolFlag) String() string {
	if f.IsSet {
		return fmt.Sprintf("%v", f.Value)
	}
	return "false"
}
