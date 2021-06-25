package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"gopkg.in/yaml.v2"
)

func NewCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "new <name>",
		Short: "create a new written",
		Run: func(cmd *cobra.Command, args []string) {
			stdout, err := handleNewCommand(args)
			if err != nil {
				fmt.Println(err)
			}
			fmt.Fprintf(cmd.OutOrStdout(), stdout)
		},
	}
}

func handleNewCommand(args []string) (string, error) {

	if len(args) != 1 {
		return fmt.Sprintf("Please provide a name of the new written"), nil
	}

	name := args[0]

	if err != nil {
		return "", err
	}

	return fmt.Sprintf("written %s created", name), nil
}

func init() {
	newCmd := NewCmd()
	rootCmd.AddCommand(newCmd)
}
