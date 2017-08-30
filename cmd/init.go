package cmd

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"

	"github.com/spf13/cobra"
	"github.com/voidint/gbb/config"
	"github.com/voidint/gbb/util"
)

var initCmd = &cobra.Command{
	Use:   "init",
	Short: `Help you to creating "gbb.json" step by step.`,
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {
		genConfigFile(confFile)
	},
}

func init() {
	RootCmd.AddCommand(initCmd)
}

func genConfigFile(destFilename string) (err error) {
	c := gather()

	fmt.Printf("About to write to %s:\n\n", destFilename)
	enc := json.NewEncoder(os.Stdout)
	enc.SetIndent("", "    ")
	if err = enc.Encode(c); err != nil {
		return err
	}
	fmt.Printf("\nIs this ok?[y/n] ")

	if ok, _ := util.Scanln(); strings.ToLower(ok) == "y" {
		return config.Save(c, destFilename)
	}
	return nil
}

func gather() (c *config.Config) {
	fmt.Printf(`This utility will walk you through creating a gbb.json file.
It only covers the most common items, and tries to guess sensible defaults.

Press ^C at any time to quit.
`)

	c = new(config.Config)
	// required
	c.Version = Version
	c.Tool = gatherOne("tool", "go install")

	fmt.Print("Do you want to continue?[y/n] ")
	if sContinue, _ := util.Scanln(); strings.ToLower(sContinue) == "n" {
		return c
	}

	// optional
	c.Importpath = gatherOne("importpath", "main")
	for {
		c.Variables = append(c.Variables, *gatherOneVar())

		fmt.Print("Do you want to continue?[y/n] ")
		if sContinue, _ := util.Scanln(); strings.ToLower(sContinue) == "n" {
			break
		}
	}
	return c
}

func gatherOneVar() (v *config.Variable) {
	return &config.Variable{
		Variable: gatherOne("variable", ""),
		Value:    gatherOne("value", ""),
	}
}

func gatherOne(prompt, defaultVal string) (input string) {
	for {
		if defaultVal != "" {
			fmt.Printf("%s: (%s) ", prompt, defaultVal)
		} else {
			fmt.Printf("%s: ", prompt)
		}
		if input, _ = util.Scanln(); input == "" {
			if defaultVal == "" {
				continue
			}
			return defaultVal
		}
		return input
	}
}
