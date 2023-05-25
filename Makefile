NAME ?= invader
INSTALL_DIR ?= /usr/local/bin

go_files := $(shell find . -iname '*.go')
bin := $(INSTALL_DIR)/$(NAME)

# build
all install: $(bin)
re: clean install

# clean
clean:; rm -f $(bin)

# test
test:; go test -v ./...

# lint
lint:;	go vet -v ./...

.PHONY: all install clean re test

###

$(bin): Makefile $(go_files)
	go build -v -o "$(bin)" ./cmd/invader/...
	@chmod +x "$@"
