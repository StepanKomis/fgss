// SPDX-License-Identifier: Apache-2.0

// Package version holds the build version.
package version

// Version is overridden at build time:
//
//	-ldflags "-X github.com/StepanKomis/fgss/src/internal/version.Version=1.2.3"
var Version = "dev"
