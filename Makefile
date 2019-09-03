test:
	go vet ./...
	go test -race `go list ./... | grep -v /vendor/`

test-verbose:
	go test -race -v `go list ./... | grep -v /vendor/`

fmt:
	godep go fmt

generate:
	gomock_generator

graph:
	go run visualize/main.go | dot -Tps -o state_machine.png
