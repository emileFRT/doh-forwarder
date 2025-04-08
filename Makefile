BIN := doh-forwarder
PREFIX := /usr/local
BUILD_FLAGS := -v -ldflags="-s -w" -trimpath 

.PHONY: build install clean

build:
	@CGO_ENABLED=0 go build $(BUILD_FLAGS) -o $(BIN)

install: build
	@install -Dm755 $(BIN) $(PREFIX)/bin/$(BIN)

clean:
	@rm -f $(BIN)