package cmd

import (
	"fmt"
	"os"
	"os/exec"
	"strings"

	"github.com/spf13/cobra"
	"github.com/voidint/gbb/util"
)

var cleanCmd = &cobra.Command{
	Use:   "clean",
	Short: `Run 'go clean' in the current directory and its subdirectories.`,
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {
		if err := rclean(wd, &cleanOpts, &gopts); err != nil {
			fmt.Fprintln(os.Stderr, err.Error())
			os.Exit(-1)
		}
	},
}

// CleanOptions clean sub-command options
type CleanOptions struct {
	I bool // go clean -i
	N bool // go clean -n
	R bool // go clean -r
	X bool // go clean -x
}

var cleanOpts CleanOptions

func init() {
	RootCmd.AddCommand(cleanCmd)
	cleanCmd.Flags().BoolVarP(&cleanOpts.I, "installed", "i", false, `To remove the corresponding installed archive or binary (what 'go install' would create).`)
	cleanCmd.Flags().BoolVarP(&cleanOpts.N, "not-execute", "n", false, `To print the remove commands it would execute, but not run them.`)
	cleanCmd.Flags().BoolVarP(&cleanOpts.R, "recursive", "r", false, `To be applied recursively to all the dependencies of the packages named by the import paths.`)
	cleanCmd.Flags().BoolVarP(&cleanOpts.X, "execute", "x", false, `To print remove commands as it executes them.`)
}

// rclean 在指定目录及其子目录下寻找main package目录并在其中执行go clean
func rclean(rootDir string, opts *CleanOptions, gopts *GlobalOptions) (err error) {
	defer util.Chdir(rootDir, gopts.Debug)

	mainPaths, err := util.WalkPkgsFunc(rootDir, util.IsMainPkg)
	if err != nil {
		return err
	}
	for i := range mainPaths {
		_ = clean(mainPaths[i], opts, gopts)
	}

	return nil
}

// clean 切换到指定目录后执行`go clean`命令
func clean(dir string, opts *CleanOptions, gopts *GlobalOptions) (err error) {
	if err = util.Chdir(dir, gopts.Debug); err != nil {
		return err
	}
	cmdArgs := []string{"go", "clean"}
	if opts != nil && opts.I {
		cmdArgs = append(cmdArgs, "-i")
	}
	if opts != nil && opts.N {
		cmdArgs = append(cmdArgs, "-n")
	}
	if opts != nil && opts.R {
		cmdArgs = append(cmdArgs, "-r")
	}
	if opts != nil && opts.X {
		cmdArgs = append(cmdArgs, "-x")
	}
	if gopts != nil && gopts.Debug {
		fmt.Printf("==> %s\n", strings.Join(cmdArgs, " "))
	}
	cmd := exec.Command(cmdArgs[0], cmdArgs[1:]...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}
