GO ?= go
GOOS ?= $(shell $(GO) env GOOS)
GOARCH ?= $(shell $(GO) env GOARCH)
MODULE_NAME ?= $(shell head -n1 go.mod | cut -f 2 -d ' ')

all: build

# Download dependencies
.PHONY: get
get:
	go mod download

# Build
#  Build for ${GOOS} and ${GOARCH}.
#  Output to `./.build/${GOOS}-${GOARCH}/*`.
.PHONY: build
build: get
	mkdir -p .build/$(GOOS)-$(GOARCH)/
	GOOS=$(GOOS) GOARCH=$(GOARCH) $(GO) build
	if [ $(GOOS) = "windows" ] ; then \
		mv ./$(MODULE_NAME).exe ./.build/$(GOOS)-$(GOARCH)/ ; \
	else \
		mv ./$(MODULE_NAME) ./.build/$(GOOS)-$(GOARCH)/ ; \
	fi ; \

# Package
#   Create an archive file containing the required files by specified OS and CPU architecture.
#   Output to `./.packages/${TAG}/${MODULE_NAME}-${TAG}.${GOOS}-${GOARCH}`
TAG ?= $(shell git tag | tail -n1)
.PHONY: package
package:
	mkdir -p ./.packages/$(TAG)/$(MODULE_NAME)-$(TAG).$(GOOS)-$(GOARCH)
	cp -r LICENSE README.md \
		./.packages/$(TAG)/$(MODULE_NAME)-$(TAG).$(GOOS)-$(GOARCH)
	if [ $(GOOS) = "windows" ] ; then \
		cp ./.build/$(GOOS)-$(GOARCH)/$(MODULE_NAME).exe ./.packages/$(TAG)/$(MODULE_NAME)-$(TAG).$(GOOS)-$(GOARCH) ; \
	else \
		cp ./.build/$(GOOS)-$(GOARCH)/$(MODULE_NAME) ./.packages/$(TAG)/$(MODULE_NAME)-$(TAG).$(GOOS)-$(GOARCH) ; \
	fi
	cd ./.packages/$(TAG) ; \
	if [ $(GOOS) = "windows" ] ; then \
		zip -r $(MODULE_NAME)-$(TAG).$(GOOS)-$(GOARCH).zip ./$(MODULE_NAME)-$(TAG).$(GOOS)-$(GOARCH) ; \
	else \
		tar cvf $(MODULE_NAME)-$(TAG).$(GOOS)-$(GOARCH).tar.gz ./$(MODULE_NAME)-$(TAG).$(GOOS)-$(GOARCH) ; \
	fi ; \
	rm -r ./$(MODULE_NAME)-$(TAG).$(GOOS)-$(GOARCH)

# Package for each OS and CPU
#   Create an archive file containing the required files per OS and CPU architecture.
#   Output to `./.packages/${TAG}/${MODULE_NAME}-${TAG}.${GOOS}-${GOARCH}`
.PHONY: package-all
package-all: get
	$(GO) tool dist list | grep 'darwin\|freebsd\|illumos\|linux\|netbsd\|openbsd\|windows' | while read line ; \
	do \
		printf GOOS= > ./.build.env ; \
		echo $$line | cut -f 1 -d "/" >> ./.build.env ; \
		printf GOARCH= >> ./.build.env ; \
		echo $$line | cut -f 2 -d "/" >> ./.build.env ; \
		. ./.build.env ; \
		make build GOOS=$$GOOS GOARCH=$$GOARCH ; \
		make package GOOS=$$GOOS GOARCH=$$GOARCH ; \
	done
	rm ./.build.env

# Claen
#   Remove artifacts in `.build/` and `.packages/`.
.PHONY: clean
clean:
	-rm -r ./.build ./.packages ./.build.env > /dev/null | true
