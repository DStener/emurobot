.DEFAULT_GOAL := build

.PHONY: build test fmt vet clean run

build: vet
	go build -o build/serialrec ./cli/serialrec

run: build 
	build/serialrec