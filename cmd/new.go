package cmd

import (
	"fmt"
	"github.com/polpettone/written/pkg"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"io/ioutil"
)

func NewCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "new <name>",
		Short: "create a new written",
		Run: func(cmd *cobra.Command, args []string) {
			stdout, err := handleNewCommand(args, cmd)
			if err != nil {
				fmt.Println(err)
			}
			fmt.Fprintf(cmd.OutOrStdout(), stdout)
		},
	}
}

func findPattern(pattern string) {
	fmt.Printf("Find pattern for name %s", pattern)
}


func handleNewCommand(args []string, cobraCommand *cobra.Command) (string, error) {

	patternName, _  := cobraCommand.Flags().GetString("pattern")
    if patternName != "" {
	    findPattern(patternName)
    }

	if len(args) != 1 {
		return fmt.Sprintf("Please provide a name of the new written"), nil
	}
	name := args[0]

	text, err := pkg.CaptureInputFromEditor("")
	if err != nil {
		return "", err
	}

	writtenDirectory := viper.GetString(WrittenDirectory)

	path := writtenDirectory + "/" + name

	err = ioutil.WriteFile(path, []byte(text), 0644)
	if err != nil {
		return "", err
	}

	return fmt.Sprintf("written %s created %s", name, writtenDirectory), nil
}

func init() {
	newCmd := NewCmd()

	newCmd.Flags().StringP(
		"pattern",
		"p",
		"",
		"pattern <name> to choose a pattern",
	)

	rootCmd.AddCommand(newCmd)
}