package main

import (
	"fmt"
	"runtime"

	"github.com/spf13/cobra"
	"honnef.co/go/tools/version"
)

func commandVersion() *cobra.Command {
	return &cobra.Command{
		Use:   "version",
		Short: "Print the version and exit",
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Printf(`commacloner Version: %s
Go Version: %s
Go OS/ARCH: %s %s
`, version.Version, runtime.Version(), runtime.GOOS, runtime.GOARCH)
		},
	}
}
