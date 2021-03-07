.PHONY: run
run: build
	./makecopy

.PHONY: build
build:
	go build .