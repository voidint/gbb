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
	// DefaultConfFile 默认配置文件路径（./gbb.json）
	DefaultConfFile = "gbb.json"
)

var (
	wd       string // 当前工作目录
	confFile string // 配置文件路径
	debug    bool   // 是否开启debug模式
)

// RootCmd represents the base command when called without any subcommands
var RootCmd = &cobra.Command{
	Use:   "gbb",
	Short: "Compile assistant",
	Long:  ``,

	Run: func(cmd *cobra.Command, args []string) {
		if confFile == DefaultConfFile {
			confFile = filepath.Join(wd, "gbb.json")
		}

		if !util.FileExist(confFile) {
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
			gt, err := util.VersionGreaterThan(Version, conf.Version)
			if err != nil {
				fmt.Fprintf(os.Stderr, err.Error())
				os.Exit(-1)
			}

			if gt { // 程序版本大于配置文件版本，重新生成配置文件。
				fmt.Printf("Warning: The gbb.json file needs to be upgraded.\n\n")
				genConfigFile(confFile)
			} else {
				// 配置文件版本大于程序版本，提醒用户升级程序。
				fmt.Printf("Warning: This program needs to be upgraded by `go get -u -v github.com/voidint/gbb`\n\n")
			}
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
	RootCmd.PersistentFlags().StringVar(&confFile, "config", DefaultConfFile, "Configuration file")
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
