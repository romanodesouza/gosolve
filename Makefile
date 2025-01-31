run:
	@go run cmd/webserver/main.go input.txt

mocks:
	go install go.uber.org/mock/mockgen@latest
	mockgen -source=internal/search/search.go -destination=internal/search/mocks/search.go
.PHONY: run