package cmd

import (
	"os"

	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "bcast",
	Short: "bcast는 오디오 스트리밍 서버로 송출할 오디오를 등록하는 CLI입니다.",
	Long:  ``,
}

func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	rootCmd.AddCommand(injectCmd)
}
