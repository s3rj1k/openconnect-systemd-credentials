## SPDX-License-Identifier: MIT
## Copyright 2024 s3rj1k.

BINARY_NAME := oc-otp
BUILD_FLAGS := -a -trimpath -buildvcs=true -ldflags '-s -w -extldflags "-static"'

export CGO_ENABLED := 0

.PHONY: all

all: build

build:
	@go build $(BUILD_FLAGS) -o $(BINARY_NAME) .
