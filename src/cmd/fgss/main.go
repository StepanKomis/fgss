// SPDX-License-Identifier: Apache-2.0

package main

import (
	"fmt"
	"os"

	"github.com/StepanKomis/fgss/src/internal/plugin"
	"github.com/StepanKomis/fgss/src/internal/remote"
	"github.com/StepanKomis/fgss/src/internal/version"
)

func usage() {
	fmt.Fprintf(os.Stderr, `Usage:
  fgss run <host> <script|path>    run a script on user@host or an ssh alias
                                   (names resolve to %s/<name>[.sh])
  fgss version
`, plugin.Dir())
}

func main() {
	if len(os.Args) < 2 {
		usage()
		os.Exit(2)
	}

	switch os.Args[1] {
	case "version":
		fmt.Println("fgss", version.Version)
	case "run":
		if len(os.Args) != 4 {
			usage()
			os.Exit(2)
		}
		if err := remote.New().Run(os.Args[2], plugin.Resolve(os.Args[3])); err != nil {
			fmt.Fprintln(os.Stderr, "error:", err)
			os.Exit(1)
		}
	case "-h", "--help", "help":
		usage()
	default:
		fmt.Fprintf(os.Stderr, "fgss: unknown command %q\n\n", os.Args[1])
		usage()
		os.Exit(2)
	}
}
