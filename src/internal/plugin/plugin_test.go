// SPDX-License-Identifier: Apache-2.0

package plugin

import (
	"os"
	"path/filepath"
	"testing"
)

func TestResolve(t *testing.T) {
	dir := t.TempDir()
	t.Setenv("FGSS_SCRIPTS_DIR", dir)

	for _, name := range []string{"health.sh", "probe", "both", "both.sh"} {
		if err := os.WriteFile(filepath.Join(dir, name), []byte("x"), 0o755); err != nil {
			t.Fatal(err)
		}
	}
	if err := os.Mkdir(filepath.Join(dir, "adir"), 0o755); err != nil {
		t.Fatal(err)
	}

	cases := map[string]string{
		"health":  filepath.Join(dir, "health.sh"), // script by name
		"probe":   filepath.Join(dir, "probe"),     // binary by exact name
		"both":    filepath.Join(dir, "both"),      // exact name wins over .sh
		"adir":    "adir",                          // directories are not plugins
		"missing": "missing",                       // falls back to the argument
	}
	for arg, want := range cases {
		if got := Resolve(arg); got != want {
			t.Errorf("Resolve(%q) = %q, want %q", arg, got, want)
		}
	}

	// An existing path is used as-is.
	p := filepath.Join(dir, "probe")
	if got := Resolve(p); got != p {
		t.Errorf("Resolve(%q) = %q, want unchanged", p, got)
	}
}
