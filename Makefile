# Makefile for golang library project

# Copyright (c) 2019-2020 Vasiliy Vasilyuk. All rights reserved.
# Use of this source code is governed by a BSD-style
# license that can be found in the LICENSE file.

# Based on https://git.io/fjkGc

# NAMESPACE the full path to the main package is use in the
# imports tool to format imports correctly.
NAMESPACE = github.com/xorcare/pointer

# COVER_FILE the name of the file recommended in the standard
# documentation go test -cover and used codecov.io
# to check code coverage.
COVER_FILE ?= coverage.out

VCS_VERSION ?= $(shell git describe --dirty --long --always --tags || echo 'undefined')

# AT addition to commands to hide unnecessary command output.
AT ?= @

build: ## Build the project binary
	$(AT)go build -ldflags "-X main.Version=$(VCS_VERSION)" ./cmd/miflib

run: ## Build the project binary
	$(AT)go run -ldflags "-X main.Version=$(VCS_VERSION)" ./cmd/miflib

check: static test build ## Check project with static checks and unit tests

help: ## Print this help
	$(AT)grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | sort | \
	awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-30s\033[0m %s\n", $$1, $$2}'

imports: tools ## Check and fix import section by import rules
	$(AT)test -z $$(for d in $$(go list -f {{.Dir}} ./...); \
	do goimports -e -l -local $(NAMESPACE) -w $$d/*.go; done)

install: ## Install the project binary
	$(AT)go install -ldflags "-X main.Version=$(VCS_VERSION)" ./cmd/miflib

lint: tools ## Check the project with lint
	$(AT)golint -set_exit_status ./...

static: imports vet lint ## Run static checks (lint, imports, vet, etc.) all over the project

test: ## Run unit tests
	$(AT)go test ./... -count=1 -race
	$(AT)go test ./... -count=1 -coverprofile=$(COVER_FILE) -covermode=atomic $d
	$(AT)go tool cover -func=$(COVER_FILE) | grep ^total

tools: ## Install all needed tools, e.g., for static checks
	$(AT)echo Installing tools from tools.go
	$(AT)cat tools.go | grep _ | awk -F'"' '{print $$2}' | xargs -tI % go install %

vet: ## Check the project with vet
	$(AT)go vet ./...

.PHONY: build check help imports install lint static test tools toolsup vet run
.DEFAULT_GOAL := build