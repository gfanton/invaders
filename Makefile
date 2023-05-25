NAME ?= invader
INSTALL_DIR ?= /usr/local/bin

go_files := $(shell find . -iname '*.go')
bin := $(INSTALL_DIR)/$(NAME)

all install: $(bin)
clean:; rm -f $(bin)
re: clean install

.PHONY: all install clean

###

$(bin): Makefile $(go_files)
	go build -v -o "$(bin)" ./cmd/invader/...
	@chmod +x "$@"
