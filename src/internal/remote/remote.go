// SPDX-License-Identifier: Apache-2.0

// Package remote copies a file to a host and executes it there, using the
// system ssh and scp so ~/.ssh/config is honoured.
package remote

import (
	"fmt"
	"io"
	"os"
	"os/exec"
	"path/filepath"
)

// Runner runs files on remote hosts. The zero value is not usable; use New.
type Runner struct {
	SCP, SSH string // executables, looked up on PATH
	Stdin    io.Reader
	Stdout   io.Writer
	Stderr   io.Writer
}

// New returns a Runner wired to the process's standard streams, so ssh
// prompts (passphrase, host key confirmation) and remote output are visible.
func New() *Runner {
	return &Runner{SCP: "scp", SSH: "ssh", Stdin: os.Stdin, Stdout: os.Stdout, Stderr: os.Stderr}
}

func (r *Runner) run(name string, args ...string) error {
	cmd := exec.Command(name, args...)
	cmd.Stdin = r.Stdin
	cmd.Stdout = r.Stdout
	cmd.Stderr = r.Stderr
	return cmd.Run()
}

// Run copies localFile to the remote host and executes it, then removes it.
// target is "user@host" or an alias from ~/.ssh/config.
func (r *Runner) Run(target, localFile string) error {
	remotePath := "/tmp/" + filepath.Base(localFile)

	// 1. Upload
	if err := r.run(r.SCP, "-q", localFile, target+":"+remotePath); err != nil {
		return fmt.Errorf("scp: %w", err)
	}

	// 2. Execute, then clean up
	remoteCmd := fmt.Sprintf("chmod +x %[1]s && %[1]s; rc=$?; rm -f %[1]s; exit $rc", remotePath)
	if err := r.run(r.SSH, target, remoteCmd); err != nil {
		return fmt.Errorf("ssh: %w", err)
	}
	return nil
}
