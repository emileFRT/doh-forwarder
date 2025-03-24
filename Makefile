BIN := doh-forwarder
PREFIX := /usr/local

.PHONY: build install clean

build:
	CGO_ENABLED=0 go build -trimpath -o $(BIN)

install: build
	install -Dm755 $(BIN) $(PREFIX)/bin/$(BIN)

clean:
	rm -f $(BIN)