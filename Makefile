BINARY=forge
GOBIN=/usr/lib/go/bin/go

.PHONY: build test bench clean install lint

build:
	$(GOBIN) build -ldflags="-s -w" -o $(BINARY) .

test:
	$(GOBIN) test ./...

bench:
	$(GOBIN) test -bench=. ./benchmarks/

clean:
	rm -f $(BINARY)

install:
	$(GOBIN) install .

lint:
	gofmt -l ./...
