SOURCES = \
  $(wildcard *.go) \
  $(wildcard main/*.go) \
  $(wildcard colors/*.go) \
  $(wildcard command/*.go) \
  $(wildcard flyterm/*.go) \
  $(wildcard periodic/*.go)

PRODUCT = _bin/attest
CHECKPOINTS = _get.ok
ARTIFACTS = $(PRODUCT) $(CHECKPOINTS)


.PHONY: all clean test

all: $(PRODUCT)
	@:

clean:
	rm -f $(ARTIFACTS)

test: _get.ok
	go vet ./...
	go test ./...

_get.ok:
	go get ./...
	@touch $@

$(PRODUCT): $(SOURCES) _get.ok
	go build -o $@ ./main
