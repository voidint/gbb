package cmd

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/spf13/cobra"
	"github.com/voidint/gbb/config"
)

var initCmd = &cobra.Command{
	Use:   "init",
	Short: "Help you to creating gbb.json step by step.",
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {
		genConfigFile(confFile)
	},
}

func init() {
	RootCmd.AddCommand(initCmd)
}

func genConfigFile(destFilename string) {
	c := gather()

	b, _ := json.MarshalIndent(c, "", "    ")
	fmt.Printf("About to write to %s:\n\n%s\n\nIs this ok?[y/n] ", destFilename, string(b))
	var ok string
	fmt.Scanln(&ok)
	if ok = strings.ToLower(ok); ok == "y" {
		config.Save(c, confFile)
	}
}

func gather() (c *config.Config) {
	fmt.Println(`This utility will walk you through creating a gbb.json file.
It only covers the most common items, and tries to guess sensible defaults.`)
	fmt.Printf("\nPress ^C at any time to quit.\n")

	c = new(config.Config)
	// required
	c.Version = Version
	c.Tool = gatherOne("tool", "go_install")

	var sContinue string
	fmt.Print("Do you want to continue?[y/n] ")
	fmt.Scanln(&sContinue)
	if sContinue = strings.ToLower(sContinue); sContinue == "n" {
		return c
	}

	// optional
	c.Importpath = gatherOne("importpath", "main")
	for {
		c.Variables = append(c.Variables, *gatherOneVar())

		fmt.Print("Do you want to continue?[y/n] ")
		fmt.Scanln(&sContinue)
		if sContinue = strings.ToLower(sContinue); sContinue == "n" {
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
		fmt.Scanln(&input) // TODO bug: 无法获取到包含空格的全部输入，如go build
		if input = strings.TrimSpace(input); input == "" {
			if defaultVal == "" {
				continue
			}
			return strings.Replace(defaultVal, "_", " ", -1)
		}
		return strings.Replace(input, "_", " ", -1) // TODO 临时举措，还原实际的空格。如，go_build ==> go build
	}
}
