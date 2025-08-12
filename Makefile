## default to coverage
all: coverage

## run test
test:
	go test -v ./...

## generate coverage + html
htmlcoverage: make-coverage-path
	go test -coverprofile=coverage/coverage.out ./...
	go tool cover -html=coverage/coverage.out -o coverage/coverage.html

## generate coverage + func output
coverage: make-coverage-path
	go test -coverprofile=coverage/coverage.out ./...
	go tool cover -func=coverage/coverage.out 

## run benchmark
bench:
	go test -bench=. ./...

make-coverage-path:
	@mkdir -p coverage

## clean
clean:
	rm -rf coverage

.PHONY: make-coverage-path test clean coverage bench
