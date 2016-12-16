package cmd

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"strings"

	"github.com/spf13/cobra"
	"github.com/voidint/gbb/config"
	"github.com/voidint/gbb/util"
	"github.com/voidint/gbb/variable"
)

var confFile string
var debug bool

// RootCmd represents the base command when called without any subcommands
var RootCmd = &cobra.Command{
	Use:   "gbb",
	Short: "",
	Long:  ``,

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
	RootCmd.PersistentFlags().BoolVar(&debug, "debug", false, "print detail")
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

	cmdArgs := strings.Fields(conf.Tool) // go install ==> []string{"go", "install"}

	var buf bytes.Buffer
	for i := range conf.Variables {
		varName := conf.Variables[i].Variable
		varExpr := conf.Variables[i].Value

		val, err := variable.Eval(varExpr)
		if err != nil {
			return err
		}
		buf.WriteString(fmt.Sprintf(`-X "%s.%s=%s"`, conf.Package, varName, val))
		if i < len(conf.Variables)-1 {
			buf.WriteByte(' ')
		}
	}
	cmdArgs = append(cmdArgs, "-ldflags", buf.String())

	if debug {
		fmt.Print("==> ", cmdArgs[0])
		args := cmdArgs[1:]
		for i := range args {
			if i < len(args)-1 {
				fmt.Print(" ", args[i])
			} else {
				fmt.Println(" ", "'"+args[i]+"'")
			}
		}
	}

	cmd := exec.Command(cmdArgs[0], cmdArgs[1:]...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}
