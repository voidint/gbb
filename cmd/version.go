package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/voidint/gbb/build"
)

var (
	// Version 版本号
	Version = "v0.0.1"
)

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Print version information",
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Printf("gbb version %s\n", Version)
		if build.Date != "" {
			fmt.Printf("date: %s\n", build.Date)
		}
		if build.Commit != "" {
			fmt.Printf("commit: %s\n", build.Commit)
		}
	},
}

func init() {
	RootCmd.AddCommand(versionCmd)
}
