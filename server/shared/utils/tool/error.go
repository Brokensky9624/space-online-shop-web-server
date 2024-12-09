package tool

import (
	"fmt"
)

func PrefixError(a string, b error) error {
	if b == nil {
		return b
	}
	if a == "" {
		a = "error"
	}
	return fmt.Errorf("%s, err: %s", a, b)
}
