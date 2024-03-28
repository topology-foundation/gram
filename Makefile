default: gram

.PHONY: gram
gram:
	go build -o ./bin/gram ./cmd/main.go
	@echo "Finished building. Run \"./bin/gram\" to launch gram."

.PHONY: clean
clean:
	rm -rf bin

.PHONY: build-proto
build-proto:
	buf generate proto