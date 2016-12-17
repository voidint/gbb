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

func gather() *config.Config {
	fmt.Println(`This utility will walk you through creating a gbb.json file.
It only covers the most common items, and tries to guess sensible defaults.`)
	fmt.Printf("\nPress ^C at any time to quit.\n")

	c := config.Config{
		Version: gatherOne("version", "0.0.1"),
		Tool:    gatherOne("tool", "go install"),
		Package: gatherOne("package", "main"),
	}
	for {
		c.Variables = append(c.Variables, *gatherOneVar())

		var sContinue string
		fmt.Print("Do you want to continue?[y/n] ")
		fmt.Scanln(&sContinue)
		if sContinue = strings.ToLower(sContinue); sContinue == "n" {
			break
		}
	}
	return &c
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
		fmt.Scanln(&input)
		if input = strings.TrimSpace(input); input == "" {
			if defaultVal == "" {
				continue
			}
			return defaultVal
		}
		return input
	}
}
