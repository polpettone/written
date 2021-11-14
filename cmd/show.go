package cmd

import (
	"fmt"
	"github.com/polpettone/written/cmd/ui"
	"github.com/spf13/cobra"
)

func ShowCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "show",
		Short: "shows all writtens",
		Run: func(cmd *cobra.Command, args []string) {
			stdout, err := handleShowCommand()
			if err != nil {
				fmt.Println(err)
			}
			fmt.Fprintf(cmd.OutOrStdout(), stdout)
		},
	}
}


func handleShowCommand() (string, error) {
	ui.MainView()
	return "", nil
}

func init() {
	showCmd := ShowCmd()
	rootCmd.AddCommand(showCmd)
}


