TARGET = wiki
PREFIX ?= /usr/local
MANPREFIX ?= "$(PREFIX)/share/man/man1"
VERSION="1.2.0"
BUILD_DIR="build"
BUILD_TARGET="$(BUILD_DIR)/$(TARGET)"
DIST="$(TARGET)-$(VERSION)"
DIST_TARGET="$(DIST).tar.gz"
DIST_SRC="$(TARGET)-src-$(VERSION)"
DIST_SRC_TARGET="$(DIST_SRC).tar.gz"

.PHONY: default all install uninstall dist dist-src clean

default: $(BUILD_TARGET)
all: clean test dist dist-src

$(BUILD_TARGET): 
	go get github.com/mattn/go-colorable
	go build -o $(BUILD_TARGET) cmd/wiki/*.go

clean:
	-rm -r $(BUILD_DIR)
	-rm -r $(DIST)
	-rm -r $(DIST_TARGET)
	-rm -r $(DIST_SRC)
	-rm -r $(DIST_SRC_TARGET)

test: $(BUILD_TARGET)
	go test -cover
	./integration-tests.sh

install:
	go install github.com/walle/wiki/cmd/wiki
	install _doc/wiki.1 $(MANPREFIX)

uninstall:
	which $(TARGET) | xargs rm
	-rm -f $(MANPREFIX)/wiki.1

dist: $(BUILD_TARGET)
	mkdir $(DIST)
	cp -r $(BUILD_TARGET) LICENSE README.md _doc $(DIST)
	tar cfz $(DIST_TARGET) $(DIST)	

dist-src:
	mkdir $(DIST_SRC)
	cp -r _doc cmd LICENSE README.md *.go integration-tests.sh Makefile $(DIST_SRC)
	tar cfz $(DIST_SRC_TARGET) $(DIST_SRC)
