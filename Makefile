BIN := doh-forwarder
PREFIX := /usr/local

.PHONY: build install clean

build:
	@echo "Building..."
	@go build -v -ldflags="-s -w" -trimpath -o $(BIN)

install: build
	@echo "Installing..."
	@install -Dm755 $(BIN) $(PREFIX)/bin/$(BIN)

clean:
	@echo "Cleaning..."
	@rm -f $(BIN)