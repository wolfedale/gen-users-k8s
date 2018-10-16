# ***********************************************************************************
# * seerene(tm) - A framework for analyzing and visualizing complex software systems.
# * Copyright (C) 2005 - 2018 for all source codes:
# * seerene(tm) GmbH, Potsdam, Germany
# ***********************************************************************************
# Go parameters
GOCMD=go
GOBUILD=$(GOCMD) build
GOCLEAN=$(GOCMD) clean
GOGET=$(GOCMD) get
BINARY=gen

all: deps build clean build-linux-i386 build-linux-amd64 build-darwin-i386 build-darwin-amd64

build:
			$(GOBUILD) -o $(BINARY) .
clean:
			$(GOCLEAN)
			rm -f $(BINARY)
run:
			$(GOBUILD) -o $(BINARY_NAME) -v ./...
			./$(BINARY_NAME)
deps:
			$(GOGET) gopkg.in/yaml.v2

clean-all:
			$(GOCLEAN)
			rm -fr releases

# Cross compilation
build-linux-i386:
	${BUILD_DIR} GOOS=linux GOARCH=386 go build -o releases/${BINARY}-linux-i386/${BINARY}-linux-i386 .

build-linux-amd64:
	${BUILD_DIR} GOOS=linux GOARCH=amd64 go build -o releases/${BINARY}-linux-amd64/${BINARY}-linux-amd64 .

build-darwin-i386:
	${BUILD_DIR} GOOS=darwin GOARCH=386 go build -o releases/${BINARY}-darwin-i386/${BINARY}-darwin-i386 .

build-darwin-amd64:
	${BUILD_DIR} GOOS=darwin GOARCH=amd64 go build -o releases/${BINARY}-darwin-amd64/${BINARY}-darwin-amd64 .
