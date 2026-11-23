// Copyright 2026 StepanKomis
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package main

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
)

// version is overridden at build time: -ldflags "-X main.version=1.2.3".
var version = "dev"

// defaultScriptsDir is where packages install the bundled scripts.
// Override with GSS_SCRIPTS_DIR (useful for development).
const defaultScriptsDir = "/etc/gss/scripts"

func scriptsDir() string {
	if d := os.Getenv("GSS_SCRIPTS_DIR"); d != "" {
		return d
	}
	return defaultScriptsDir
}

// resolveScript returns arg as-is when it is an existing file; otherwise it
// looks for <scriptsDir>/<arg>.sh so installed scripts can be used by name.
func resolveScript(arg string) string {
	if _, err := os.Stat(arg); err == nil {
		return arg
	}
	p := filepath.Join(scriptsDir(), arg+".sh")
	if _, err := os.Stat(p); err == nil {
		return p
	}
	return arg
}

// run executes a command, wiring its output to ours so the user sees
// ssh prompts (passphrase, host key confirmation) and remote output.
func run(name string, args ...string) error {
	cmd := exec.Command(name, args...)
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}

// RunScript copies a local script to the remote host and executes it.
// target is "user@host" or an alias from ~/.ssh/config.
func RunScript(target, localScript string) error {
	remotePath := "/tmp/" + filepath.Base(localScript)

	// 1. Upload
	if err := run("scp", "-q", localScript, target+":"+remotePath); err != nil {
		return fmt.Errorf("scp: %w", err)
	}

	// 2. Execute, then clean up
	remoteCmd := fmt.Sprintf("chmod +x %[1]s && %[1]s; rc=$?; rm -f %[1]s; exit $rc", remotePath)
	if err := run("ssh", target, remoteCmd); err != nil {
		return fmt.Errorf("ssh: %w", err)
	}
	return nil
}

func usage() {
	fmt.Fprintf(os.Stderr, `Usage:
  gss run <host> <script|path>     run a script on user@host or an ssh alias
                                   (names resolve to %s/<name>.sh)
  gss version
`, scriptsDir())
}

func main() {
	if len(os.Args) < 2 {
		usage()
		os.Exit(2)
	}

	switch os.Args[1] {
	case "version":
		fmt.Println("gss", version)
	case "run":
		if len(os.Args) != 4 {
			usage()
			os.Exit(2)
		}
		if err := RunScript(os.Args[2], resolveScript(os.Args[3])); err != nil {
			fmt.Fprintln(os.Stderr, "error:", err)
			os.Exit(1)
		}
	case "-h", "--help", "help":
		usage()
	default:
		fmt.Fprintf(os.Stderr, "gss: unknown command %q\n\n", os.Args[1])
		usage()
		os.Exit(2)
	}
}
