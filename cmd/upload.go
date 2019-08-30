package cmd

import (
	"fmt"

	"github.com/hellojukay/gors/upload"
	"github.com/spf13/cobra"
)

var uFilename string
var uploadCmd = &cobra.Command{
	Use:   "upload",
	Short: "upload your recording file",
	Long:  "upload your recording file",
	Run: func(cmd *cobra.Command, args []string) {
		if uFilename == "" {
			fmt.Println(cmd.UsageString())
			return
		}
		uploader := upload.NewUploader(uFilename)
		uploader.Execute()
	},
}

func init() {
	RootCmd.AddCommand(uploadCmd)
	uploadCmd.Flags().StringVarP(&uFilename, "filename", "f", "", "the file which to save your terminal data")
}
