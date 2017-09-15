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

const (
	// DefaultConfFile default configuration file path
	DefaultConfFile = "gbb.json"
)

// RootCmd represents the base command when called without any subcommands
var RootCmd = &cobra.Command{
	Use:   "gbb",
	Short: "Compile assistant",
	Long:  ``,

	Run: func(cmd *cobra.Command, args []string) {
		if gopts.ConfigFile == DefaultConfFile {
			gopts.ConfigFile = filepath.Join(wd, "gbb.json")
		}

		if !util.FileExist(gopts.ConfigFile) {
			genConfigFile(gopts.ConfigFile)
			return
		}
		conf, err := config.Load(gopts.ConfigFile)
		if err != nil {
			fmt.Fprintln(os.Stderr, err.Error())
			os.Exit(-1)
		}
		conf.Debug = gopts.Debug
		conf.All = gopts.All

		if conf.Version != Version {
			gt, err := util.VersionGreaterThan(Version, conf.Version)
			if err != nil {
				fmt.Fprintln(os.Stderr, err.Error())
				os.Exit(-1)
			}

			if gt { // 程序版本大于配置文件版本，重新生成配置文件。
				fmt.Printf("Warning: The gbb.json file needs to be upgraded.\n\n")
				genConfigFile(gopts.ConfigFile)
			} else {
				// 配置文件版本大于程序版本，提醒用户升级程序。
				fmt.Printf("Warning: This program needs to be upgraded by `go get -u -v github.com/voidint/gbb`\n\n")
			}
			return
		}

		if err := tool.Build(conf, wd); err != nil {
			fmt.Fprintln(os.Stderr, err.Error())
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

// GlobalOptions global options
type GlobalOptions struct {
	All        bool
	Debug      bool
	ConfigFile string
}

var (
	wd    string // current work directory
	gopts GlobalOptions
)

func init() {
	cobra.OnInitialize(initConfig)
	RootCmd.PersistentFlags().BoolVarP(&gopts.All, "all", "a", false, "Act on all go packages")
	RootCmd.PersistentFlags().BoolVar(&gopts.Debug, "debug", false, "Enable debug mode")
	RootCmd.PersistentFlags().StringVar(&gopts.ConfigFile, "config", DefaultConfFile, "Configuration file")
}

// initConfig reads in config file and ENV variables if set.
func initConfig() {
	var err error
	wd, err = os.Getwd()
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(-1)
	}
}
