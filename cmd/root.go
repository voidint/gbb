package cmd

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"strings"
	"time"

	"github.com/spf13/cobra"
	"github.com/voidint/gbb/config"
	"github.com/voidint/gbb/util"
)

var confFile string
var debug bool

// RootCmd represents the base command when called without any subcommands
var RootCmd = &cobra.Command{
	Use:   "gbb",
	Short: "A brief description of your application",
	Long: `A longer description that spans multiple lines and likely contains
examples and usage of using your application. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	// Uncomment the following line if your bare application
	// has an action associated with it:
	Run: func(cmd *cobra.Command, args []string) {
		if !util.FileExist("./gbb.json") {
			collect()
			return
		}
		if err := build(); err != nil {
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

	// Here you will define your flags and configuration settings.
	// Cobra supports Persistent Flags, which, if defined here,
	// will be global for your application.

	RootCmd.PersistentFlags().BoolVar(&debug, "debug", false, "print detail")
	// Cobra also supports local flags, which will only run
	// when this action is called directly.
	// RootCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}

// initConfig reads in config file and ENV variables if set.
func initConfig() {
	wd, err := os.Getwd()
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(-1)
	}
	sep := fmt.Sprintf("%c", os.PathSeparator)
	if !strings.HasSuffix(wd, sep) {
		wd += sep
	}

	confFile = wd + "gbb.json"
}

func build() (err error) {
	conf, err := config.Load(confFile)
	if err != nil {
		return err
	}

	fields := strings.Fields(conf.Tool) // go install ==> []string{"go", "install"}
	cmdName := fields[0]

	var args []string
	if len(fields) > 1 {
		args = append(args, fields[1:]...)
	}

	var buf bytes.Buffer
	for i := range conf.Variables {
		variable := conf.Variables[i].Variable
		value := conf.Variables[i].Value

		switch value {
		case "{{.date}}":
			value = time.Now().Format(time.RFC3339)
		case "{{.gitCommit}}":
			value = "d409ffce59c67b544890b89b956fe8f8a38d6c7b"
		}
		buf.WriteString(fmt.Sprintf(`-X "%s.%s=%s"`, conf.Package, variable, value))
		if i < len(conf.Variables)-1 {
			buf.WriteByte(' ')
		}
	}
	args = append(args, "-ldflags", buf.String())

	if debug {
		fmt.Print(cmdName)
		for i := range args {
			if i < len(args)-1 {
				fmt.Print(" ", args[i])
			} else {
				fmt.Println(" ", "'"+args[i]+"'")
			}
		}

	}
	cmd := exec.Command(cmdName, args...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	return cmd.Run()
}
