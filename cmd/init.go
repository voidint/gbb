package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/voidint/gbb/config"
)

var initCmd = &cobra.Command{
	Use:   "init",
	Short: "一步步搜集信息并在当前目录下gbb.json文件",
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {
		collect()
	},
}

func init() {
	RootCmd.AddCommand(initCmd)
}

func collect() {
	fmt.Println(`This utility will walk you through creating a gbb.json file.
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
