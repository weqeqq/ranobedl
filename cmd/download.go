package cmd

import (
	"fmt"
	"os"
	"ranobedl/api/ranobelib"
	"ranobedl/cachemgr"
	"ranobedl/format"
	"ranobedl/ranobe"

	"github.com/schollz/progressbar/v3"
	"github.com/spf13/cobra"
)

type downloader struct {
	Cmd  *cobra.Command
	Args []string
}

func newDownloader(cmd *cobra.Command, args []string) *downloader {
	return &downloader{cmd, args}
}

const DownloaderUrlIndex = 0

func (self *downloader) getUrl() string {
	return self.Args[DownloaderUrlIndex]
}
func (self *downloader) getOutput() string {
	output, _ := self.Cmd.Flags().GetString("output")
	return output
}
func (self *downloader) Run() error {

	uniqueName, err := ranobelib.GetUniqueName(self.getUrl())
	if err != nil {
		return err
	}
	progressbar := progressbar.NewOptions(100,
		progressbar.OptionEnableColorCodes(true),
		progressbar.OptionSetPredictTime(true),
		progressbar.OptionSetElapsedTime(true),
		progressbar.OptionSetWidth(25),
		progressbar.OptionSetDescription("Downloading..."),
		progressbar.OptionClearOnFinish(),
		progressbar.OptionSetTheme(progressbar.Theme{
			Saucer:        "[green]=[reset]",
			SaucerHead:    "[green]>[reset]",
			SaucerPadding: " ",
			BarStart:      "[",
			BarEnd:        "]",
		}),
	)
	callback := func(current, total int) {
		progressbar.Set(int(
			float64(current) / float64(total-1) * 100,
		))
	}
	if err := ranobe.Download(cachemgr.RanobeLib, uniqueName, callback); err != nil {
		return err
	}
	if err := format.Export(cachemgr.RanobeLib, uniqueName, self.getOutput()); err != nil {
		return err
	}
	fmt.Println("Success!")
	return nil
}
func runDownloadCmd(cmd *cobra.Command, args []string) {
	if err := newDownloader(cmd, args).Run(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

var downloadCmd = &cobra.Command{
	Use:   "download",
	Short: "Download ranobe",
	Long:  "Download ranobe",
	Args:  cobra.ExactArgs(1),
	Run:   runDownloadCmd,
}

func init() {
	downloadCmd.Flags().StringP(
		"format",
		"f",
		"fb2",
		"format (fb2, epub)",
	)
	downloadCmd.Flags().StringP(
		"output",
		"o",
		"ranobe.fb2",
		"output path",
	)
}
