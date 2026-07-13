.PHONY: test test-cover

test:
	go test ./... -v

test-cover:
	go test ./... -cover
