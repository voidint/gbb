package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/voidint/gbb/config"
)

// initCmd represents the init command
var initCmd = &cobra.Command{
	Use:   "init",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: func(cmd *cobra.Command, args []string) {
		// TODO: Work your own magic here
		collect()
	},
}

func init() {
	RootCmd.AddCommand(initCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// initCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// initCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")

}

func collect() {
	fmt.Println(`This utility will walk you through creating a package.json file.
It only covers the most common items, and tries to guess sensible defaults.`)

	config.Save(&config.Config{
		Version: "0.0.1",
		Tool:    "gb",
		Package: "main",
		Variables: []config.Variable{
			{
				Variable: "Date",
				Value:    "{{.date}}",
			},
			{
				Variable: "Commit",
				Value:    "{{.gitCommit}}",
			},
		},
	}, confFile)
}
