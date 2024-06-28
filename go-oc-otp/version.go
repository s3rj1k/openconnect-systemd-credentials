// SPDX-License-Identifier: MIT
// Copyright 2024 s3rj1k.

package main

import (
	"fmt"
	"runtime/debug"
	"strings"
	"time"
)

const revisionNumberCount = 10

func GetVCSBuildInfo() string {
	buildInfo, ok := debug.ReadBuildInfo()
	if !ok {
		return "unknown"
	}

	var (
		vcsRevision []rune // vcs.revision
		vcsTime     string // vcs.time
		vcsModified string // vcs.modified
	)

	for _, el := range buildInfo.Settings {
		switch el.Key {
		case "vcs.revision":
			vcsRevision = []rune(el.Value)
		case "vcs.time":
			vcsTime = el.Value
		case "vcs.modified":
			vcsModified = el.Value
		default:
			continue
		}
	}

	var revision string

	if len(vcsRevision) <= revisionNumberCount {
		revision = string(vcsRevision)
	} else {
		revision = string(vcsRevision[:revisionNumberCount])
	}

	t, err := time.Parse(time.RFC3339, vcsTime)
	if err == nil {
		revision = fmt.Sprintf("%s-%d", revision, t.Unix())
	}

	if strings.EqualFold(vcsModified, "true") {
		revision += "-dirty"
	}

	return revision
}
