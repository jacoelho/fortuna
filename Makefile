build:
	go install -gcflags '-m'

simplify:
	gofmt  -s -w ./..

test:
	go test -v -cover -race

bench:
	go test -v -cover -race -test.bench=. -test.benchmem

vet:
	go vet