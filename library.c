// SPDX-License-Identifier: MIT
// Copyright 2024 s3rj1k.

#define _GNU_SOURCE
#include <dlfcn.h>
#include <limits.h>
#include <stdio.h>
#include <stdlib.h>
#include <string.h>
#include <unistd.h>

static char *(*original_fgets)(char *buf, int n, FILE *stream) = NULL;

char *fgets(char *buf, int n, FILE *stream)
{
    if (!original_fgets)
    {
        original_fgets = (char *(*)(char *, int, FILE *))dlsym(RTLD_NEXT, "fgets");
    }

    if (stream != stdin)
    {
        return original_fgets(buf, n, stream);
    }

    const char *command = getenv("AUTH_CMD");
    if (command == NULL)
    {
        return original_fgets(buf, n, stream);
    }

    FILE *proc = popen(command, "r");
    if (proc == NULL)
    {
        return original_fgets(buf, n, stream);
    }

    char *result = original_fgets(buf, n, proc);

    int status = pclose(proc);
    if (status == -1)
    {
        return NULL;
    }

    return result;
}

char *gets(char *buf)
{
    return fgets(buf, INT_MAX, stdin) ? buf : NULL;
}
