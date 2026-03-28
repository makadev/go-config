## default to coverage
all: coverage

## run test
test:
	go test -v .

test-race:
	go test -v -race .

## generate coverage + html
htmlcoverage:
	@mkdir -p coverage
	go test -coverprofile=coverage/coverage.out .
	go tool cover -html=coverage/coverage.out -o coverage/coverage.html

## generate coverage + func output
coverage:
	@mkdir -p coverage
	go test -coverprofile=coverage/coverage.out .
	go tool cover -func=coverage/coverage.out 

## run benchmark
bench:
	go test -bench=. .

## clean
clean:
	rm -rf coverage

.PHONY: all test clean coverage htmlcoverage bench
