package cmd

import (
	"fmt"
	"github.com/polpettone/written/cmd/config"
	"github.com/polpettone/written/pkg"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"io/ioutil"
	"os"
	"path/filepath"
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

type PatternFile struct {
	info os.FileInfo
	path string
}

func findPattern(pattern string) (string, error) {
	fmt.Printf("Find pattern for name %s \n", pattern)
	patternsDirectory := viper.GetString(config.PatternsDirectory)

	var patternFiles []PatternFile

	err := filepath.Walk(patternsDirectory,
		func(path string, info os.FileInfo, err error) error {
			patternFile := PatternFile{
				info: info,
				path: path,
			}
			patternFiles = append(patternFiles, patternFile)
			return nil
		})

	if err != nil {
		return "", err
	}

	var foundPatternFile *PatternFile
	for _, file := range patternFiles {
		if file.info.Name() == pattern {
			foundPatternFile = &file
			break
		}
	}

	if foundPatternFile != nil {
		patternContent, err := ioutil.ReadFile(foundPatternFile.path)
		if err != nil {
			return "", err
		}
		return string(patternContent), nil
	}
	return "", nil
}


func handleNewCommand(args []string, cobraCommand *cobra.Command) (string, error) {

	if len(args) != 1 {
		return fmt.Sprintf("Please provide a name of the new written"), nil
	}
	name := args[0]

	patternName, _  := cobraCommand.Flags().GetString("pattern")
	var patternContent string
    if patternName != "" {
	    content, err := findPattern(patternName)
	    if err != nil {
	    	return "", err
		}
		patternContent = content
    }

	text, err := pkg.CaptureInputFromEditor(patternContent)
	if err != nil {
		return "", err
	}

	writtenDirectory := viper.GetString(config.WrittenDirectory)

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
