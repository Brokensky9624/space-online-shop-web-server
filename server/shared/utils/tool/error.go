package tool

import "fmt"

func MergeErrors(a, b error) error {
	if a == nil {
		return b
	}
	if b == nil {
		return a
	}
	return fmt.Errorf("%s; %s", a, b)
}

func PrefixError(a string, b error) error {
	if b == nil {
		return b
	}
	if a == "" {
		a = "error"
	}
	return fmt.Errorf("%s, err: %s", a, b)
}
