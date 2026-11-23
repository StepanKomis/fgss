// SPDX-License-Identifier: Apache-2.0

// Package plugin finds the scripts and binaries fgss can run on a host.
package plugin

import (
	"os"
	"path/filepath"
)

// DefaultDir is where packages install the bundled plugins.
const DefaultDir = "/etc/fgss/scripts"

// Dir returns the plugin directory. FGSS_SCRIPTS_DIR overrides the default
// (useful for development).
func Dir() string {
	if d := os.Getenv("FGSS_SCRIPTS_DIR"); d != "" {
		return d
	}
	return DefaultDir
}

func isFile(p string) bool {
	fi, err := os.Stat(p)
	return err == nil && fi.Mode().IsRegular()
}

// Resolve returns arg as-is when it is an existing file. Otherwise it looks
// in Dir for a file named exactly <arg> (any script or binary) and then for
// <arg>.sh, so installed plugins can be used by name. If nothing matches, arg
// is returned unchanged.
func Resolve(arg string) string {
	if isFile(arg) {
		return arg
	}
	for _, name := range []string{arg, arg + ".sh"} {
		if p := filepath.Join(Dir(), name); isFile(p) {
			return p
		}
	}
	return arg
}
