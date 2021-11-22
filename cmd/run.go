package cmd

import (
	"fmt"
	"github.com/polpettone/written/cmd/ui"
	"github.com/spf13/cobra"
)

func RunCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "run",
		Short: "",
		Run: func(cmd *cobra.Command, args []string) {
			stdout, err := handleRunCommand()
			if err != nil {
				fmt.Println(err)
			}
			fmt.Fprintf(cmd.OutOrStdout(), stdout)
		},
	}
}


func handleRunCommand() (string, error) {
	ui.FlexView()
	return "", nil
}

func init() {
	runCmd := RunCmd()
	rootCmd.AddCommand(runCmd)
}


