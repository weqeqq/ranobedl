package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

func runRootCmd(_ *cobra.Command, _ []string) {
	fmt.Println("See 'ranobedl --help'")
}

var rootCmd = &cobra.Command{
	Use:   "ranobedl",
	Short: "ranobedl is a CLI tool for downloading ranobe.",
	Long:  "ranobedl is a CLI tool for downloading ranobe.",
	Run:   runRootCmd,
}

func init() {
	rootCmd.AddCommand(downloadCmd)
	rootCmd.AddCommand(clearCmd)
}
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
