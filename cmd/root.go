package cmd

import (
	"fmt"
	"os"
	"path/filepath"
	"time"

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
	Use: "gbb",
	Long: `Go project compilation assistant.
Copyright (c) 2016-2020, voidint. All rights reserved.`,

	Run: func(cmd *cobra.Command, args []string) {
		begin := time.Now()

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
			return
		}
		conf.Debug = gopts.Debug
		conf.All = gopts.All

		if conf.Version != Version {
			gt, err := util.VersionGreaterThan(Version, conf.Version)
			if err != nil {
				fmt.Fprintln(os.Stderr, err.Error())
				os.Exit(-1)
				return
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

		defer func() {
			fmt.Printf("Time Used: %.2fs\n", time.Now().Sub(begin).Seconds())
		}()

		if err := tool.Build(conf, wd); err != nil {
			fmt.Fprintln(os.Stderr, err.Error())
			os.Exit(-1)
			return
		}
	},
}

// Execute adds all child commands to the root command sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	if err := RootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(-1)
		return
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
	RootCmd.PersistentFlags().BoolVarP(&gopts.All, "all", "a", false, "build all packages")
	RootCmd.PersistentFlags().BoolVarP(&gopts.Debug, "debug", "D", false, "enable debug mode")
	RootCmd.PersistentFlags().StringVarP(&gopts.ConfigFile, "config", "c", DefaultConfFile, "configuration file")
}

// initConfig reads in config file and ENV variables if set.
func initConfig() {
	var err error
	wd, err = os.Getwd()
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(-1)
		return
	}
}
