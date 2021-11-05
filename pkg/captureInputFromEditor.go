package pkg

import (
	"io/ioutil"
	"os"
	"os/exec"
	"strings"
)

const DefaultEditor = "vim"

func CaptureInputFromEditor(content string) (string, error) {
	file, err := ioutil.TempFile(os.TempDir(), "*")
	if err != nil {
		return "", err
	}

	filename := file.Name()

	defer os.Remove(filename)

	_, err = file.WriteString(content)

	if err != nil {
		return "", err
	}

	if err = file.Close(); err != nil {
		return "", err
	}

	if err = OpenFileInEditor(filename); err != nil {
		return "", err
	}

	bytes, err := ioutil.ReadFile(filename)
	if err != nil {
		return "", err
	}

	return strings.TrimSpace(string(bytes)), nil
}

func OpenFileInEditor(filename string) error {
	editor := os.Getenv("EDITOR")
	if editor == "" {
		editor = DefaultEditor
	}

	executable, err := exec.LookPath(editor)
	if err != nil {
		return err
	}

	command := exec.Command(executable, filename)
	command.Stdin = os.Stdin
	command.Stdout = os.Stdout
	command.Stderr = os.Stderr

	return command.Run()
}
