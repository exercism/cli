.PHONY: test update

GOTEST := ${shell which gotestsum 2> /dev/null}
ifdef GOTEST
	GOTEST += --
else
	GOTEST := go test
endif

test:
	${GOTEST} -vet=all ./...
	go run ./exercism troubleshoot

update:
	go get -u -t all
	# Keep in mind that by default, tidy acts as if the
	# -compat flag were set to the version prior to the
	# one indicated by the 'go' directive in the go.mod file	
	go mod tidy
