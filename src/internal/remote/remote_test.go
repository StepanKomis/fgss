// SPDX-License-Identifier: Apache-2.0

package remote

import (
	"io"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

// stub writes an executable that appends its arguments to log and exits rc.
func stub(t *testing.T, dir, name, log string, rc string) string {
	t.Helper()
	p := filepath.Join(dir, name)
	body := "#!/bin/sh\necho \"" + name + " $*\" >> '" + log + "'\nexit " + rc + "\n"
	if err := os.WriteFile(p, []byte(body), 0o755); err != nil {
		t.Fatal(err)
	}
	return p
}

func newRunner(scp, ssh string) *Runner {
	return &Runner{SCP: scp, SSH: ssh, Stdin: strings.NewReader(""), Stdout: io.Discard, Stderr: io.Discard}
}

func TestRunUploadsThenExecutes(t *testing.T) {
	dir := t.TempDir()
	log := filepath.Join(dir, "log")
	r := newRunner(stub(t, dir, "scp", log, "0"), stub(t, dir, "ssh", log, "0"))

	if err := r.Run("alice@host", "/etc/fgss/scripts/health.sh"); err != nil {
		t.Fatal(err)
	}

	b, _ := os.ReadFile(log)
	lines := strings.Split(strings.TrimSpace(string(b)), "\n")
	if len(lines) != 2 {
		t.Fatalf("want scp then ssh, got %q", lines)
	}
	if want := "scp -q /etc/fgss/scripts/health.sh alice@host:/tmp/health.sh"; lines[0] != want {
		t.Errorf("scp call = %q, want %q", lines[0], want)
	}
	if !strings.HasPrefix(lines[1], "ssh alice@host chmod +x /tmp/health.sh") {
		t.Errorf("ssh call = %q", lines[1])
	}
}

func TestRunStopsWhenUploadFails(t *testing.T) {
	dir := t.TempDir()
	log := filepath.Join(dir, "log")
	r := newRunner(stub(t, dir, "scp", log, "1"), stub(t, dir, "ssh", log, "0"))

	err := r.Run("host", "x.sh")
	if err == nil || !strings.HasPrefix(err.Error(), "scp:") {
		t.Fatalf("want scp error, got %v", err)
	}
	if b, _ := os.ReadFile(log); strings.Contains(string(b), "ssh ") {
		t.Errorf("ssh must not run after a failed upload: %q", b)
	}
}
