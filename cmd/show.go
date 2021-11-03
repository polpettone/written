package cmd

import (
	"fmt"
	"github.com/polpettone/written/cmd/models"
	"github.com/polpettone/written/cmd/service"
	"github.com/polpettone/written/cmd/ui"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
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
	documents, _ := readDocuments()
	ui.MainView(documents)
	return "", nil
}

func readDocuments() ([]*models.Document, error) {
	WrittenDirectory := viper.GetString(WrittenDirectory)
	documents, err := service.Read(WrittenDirectory)
	if err != nil {
		return nil, err
	}
	return documents, nil
}

func init() {
	showCmd := ShowCmd()
	rootCmd.AddCommand(showCmd)
}


