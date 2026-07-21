.PHONY: all run run-file build clean

BINARY=voxstream

build:
	go build -o $(BINARY) .

run: build
	./$(BINARY) -stream

run-file: build
	./$(BINARY) -file=$(FILE)

all: run

clean:
	rm -f $(BINARY)
