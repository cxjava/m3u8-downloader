package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

var (
	Version            = "dev"
	Commit             = "none"
	RepoUrl            = "unknown"
	BuildDate          = "unknown"
	BuiltBy            = "unknown"
	BuiltWithGoVersion = "unknown"
)

// versionCmd represents the version command
var versionCmd = &cobra.Command{
	Use:     "version",
	Aliases: []string{"v"},
	Short:   "Show the version.",
	Long:    `Display the version.`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("Version: \t" + Version)
		fmt.Println("Repo URL: \t" + RepoUrl)
		fmt.Println("Commit Info: \t" + Commit)
		fmt.Println("Build Time: \t" + BuildDate)
		fmt.Println("Build By: \t" + BuiltBy)
		fmt.Println("Go Version: \t" + BuiltWithGoVersion)
	},
}

func init() {
	rootCmd.AddCommand(versionCmd)
}
