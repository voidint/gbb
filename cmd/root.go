package cmd

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/spf13/cobra"
	"github.com/voidint/gbb/config"
	"github.com/voidint/gbb/tool"
	"github.com/voidint/gbb/util"
)

var confFile, wd string
var debug bool

// RootCmd represents the base command when called without any subcommands
var RootCmd = &cobra.Command{
	Use:   "gbb",
	Short: "",
	Long:  ``,

	Run: func(cmd *cobra.Command, args []string) {
		if !util.FileExist("./gbb.json") {
			genConfigFile(confFile)
			return
		}
		conf, err := config.Load(confFile)
		if err != nil {
			fmt.Fprintf(os.Stderr, err.Error())
			os.Exit(-1)
		}
		conf.Debug = debug

		if conf.Version != Version {
			fmt.Printf("The gbb.json file needs to be upgraded.\n\n")
			genConfigFile(confFile)
			return
		}

		if err := tool.Build(conf, wd); err != nil {
			fmt.Fprintf(os.Stderr, err.Error())
			os.Exit(-1)
		}
	},
}

// Execute adds all child commands to the root command sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	if err := RootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(-1)
	}
}

func init() {
	cobra.OnInitialize(initConfig)
	RootCmd.PersistentFlags().BoolVar(&debug, "debug", false, "Enable debug mode")
}

// initConfig reads in config file and ENV variables if set.
func initConfig() {
	var err error
	wd, err = os.Getwd()
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(-1)
	}
	confFile = filepath.Join(wd, "gbb.json")
}
