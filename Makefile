.PHONY: test build lint deps gen

build:
		go build ./...

test:
		go test -v ./...

# Go lint
lint:
		golangci-lint run

gen:
		mockery --all --recursive --with-expecter 

clean:
		rm -r mocks