package main

import (
	"fmt"
)

func FetchEndpoint(args []string) string {
	if len(args) == 0 {
		return FetchEndpoints["current"]
	}

	endpoint := FetchEndpoints["exercise"]
	for _, arg := range args {
		endpoint = fmt.Sprintf("%s/%s", endpoint, arg)
	}

	return endpoint
}
