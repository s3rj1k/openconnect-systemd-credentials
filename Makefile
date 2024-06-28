## SPDX-License-Identifier: MIT
## Copyright 2024 s3rj1k.

CC = gcc
CFLAGS = -shared -fPIC
SRC = library.c
TARGET = oc-password.so

all: clean format build strip

format:
	clang-format --style=Microsoft --sort-includes -i $(SRC)

build:
	$(CC) $(CFLAGS) -o $(TARGET) $(SRC) -ldl

strip:
	strip --strip-unneeded $(TARGET)

clean:
	rm -f $(TARGET)
