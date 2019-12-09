package cmd

import (
	"os"
	"path/filepath"

	"github.com/mritd/mmh/mmh"
	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:              "mmh",
	Short:            "a simple multi-server ssh tool",
	Long:             "a simple multi-server ssh tool.",
	TraverseChildren: true,
	Run:              func(cmd *cobra.Command, args []string) { mmh.InteractiveLogin() },
}

func Execute() {

	runCmd := rootCmd
	mmh.Aliases = findAllAliases(rootCmd)

	subCmd, _, err := rootCmd.Find([]string{filepath.Base(os.Args[0])})
	if err == nil && subCmd.Name() != rootCmd.Name() {
		runCmd = subCmd
		rootCmd.SetArgs(append([]string{subCmd.Name()}, os.Args[1:]...))
	}

	if runCmd.Name() != "install" && runCmd.Name() != "uninstall" {
		mmh.LoadConfig()
	}

	if err := runCmd.Execute(); err != nil {
		mmh.Exit(err.Error(), -1)
	}
}

func findAllAliases(cmd *cobra.Command) []string {
	var aliases []string
	if cmd.HasSubCommands() {
		cmds := cmd.Commands()
		for _, c := range cmds {
			if len(c.Aliases) > 0 {
				aliases = append(aliases, c.Aliases...)
			}
			if c.HasSubCommands() {
				as := findAllAliases(c)
				if len(as) > 0 {
					aliases = append(aliases, as...)
				}
			}
		}
	} else {
		if len(cmd.Aliases) > 0 {
			aliases = append(aliases, cmd.Aliases...)
		}
	}

	return aliases
}
