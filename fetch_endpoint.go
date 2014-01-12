package main

import (
	"fmt"
)

func FetchEndpoint(args []string) string {
	if len(args) == 0 {
		return FetchEndpoints["current"]
	}
	return fmt.Sprintf("%s/%s/%s", FetchEndpoints["exercise"], args[0], args[1])
}
