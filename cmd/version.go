package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/voidint/gbb/build"
)

var (
	// Version 版本号
	Version = "0.4.0"
)

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Print version information",
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println(build.Version(fmt.Sprintf("gbb version %s", Version)))
	},
}

func init() {
	RootCmd.AddCommand(versionCmd)
}
