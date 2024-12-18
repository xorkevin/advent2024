.PHONY: run build bench go run-go build-go rs run-rs build-rs

BENCHARGS=--warmup=8 --shell=none --time-unit=millisecond
BENCH=hyperfine
BIN=$(notdir $(CURDIR))

GOBIN=./bin/$(BIN)
RSBIN=../target/release/$(BIN)

run:
	go run .

build: build-go build-rs

bench: build
	$(BENCH) $(BENCHARGS) -L bin $(GOBIN),$(RSBIN) '{bin}'

go: build-go run-go

run-go:
	$(GOBIN)

build-go:
	go build -o $(GOBIN)

rs: build-rs run-rs

run-rs:
	$(RSBIN)

build-rs:
	cargo build --release
