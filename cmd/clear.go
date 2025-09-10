package cmd

import (
	"fmt"
	"os"
	"ranobedl/cachemgr"

	"github.com/spf13/cobra"
)

type clearer struct {
	Cmd  *cobra.Command
	Args []string
}

func newClearer(cmd *cobra.Command, args []string) *clearer {
	return &clearer{cmd, args}
}

func (self *clearer) Run() error {
	return cachemgr.ClearCache()
}

func runClearCmd(cmd *cobra.Command, args []string) {
	if err := newClearer(cmd, args).Run(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

var clearCmd = &cobra.Command{
	Use:   "clear",
	Short: "Clear cache",
	Long:  "Clear cache",
	Run:   runClearCmd,
}
