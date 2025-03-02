package file

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

const userFilesDir = "userFiles"

var forbbiden_extensions []string = []string{".go", ".exe", ".dll", ".sh", ".md"}

func doesFileExists(file_name string) bool {
	_, error := os.Stat(file_name)
	return !os.IsNotExist(error)

}

func readFromFile(fileName string) (string, error) {
	exists := doesFileExists(fileName)
	if !exists {
		return "", ErrFileDoesNotExist
	}

	data, err := os.ReadFile(fileName)
	if err != nil {
		return "", err
	}

	return string(data), nil
}

func WriteOutput(outputFile string, fileContent string) error {
	file, err := os.OpenFile(outputFile, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0644)
	if err != nil {
		return ErrCannotCreateFile
	}
	defer file.Close()

	_, err = file.WriteString(fileContent)
	if err != nil {
		return ErrCouldntWriteFile
	}

	return nil
}

func isForbidden(file_path string) bool {
	extension := strings.ToLower(filepath.Ext(file_path))
	for _, ext := range forbbiden_extensions {
		if extension == ext {
			return true
		}
	}
	return false
}

func CheckArgument(words []string) ([]string, string, error) {
	if words[0] == "prompt" || words[0] == "touch" || words[0] == "truncate" || words[0] == "rm" {
		return words, "", nil
	}
	var arg string
	var err error

	for i, word := range words[1:] {

		if strings.Contains(word, "\"") || word[0] == '>' || word[0] == '-' {
			continue
		}
		words = RemoveAtIndex(words, i+1)
		word = strings.ReplaceAll(word, "<", "")
		arg, err = readFromFile(word)
		arg = `"` + arg + `"`
		if err != nil {
			return words, "", ErrFileDoesNotExist
		}

	}
	return words, arg, nil
}

func RemoveAtIndex(s []string, i int) []string {
	if i < 0 || i >= len(s) {
		return s
	}
	return append(s[:i], s[i+1:]...)
}

func EnsureUserFilesDir() error {
	fmt.Println("here")
	absPath, err := filepath.Abs(userFilesDir)
	if err != nil {
		return err
	}
	// MkdirAll creates the directory along with any missing parents.
	// It does nothing if the directory already exists.
	if err := os.MkdirAll(absPath, 0700); err != nil {
		emkd := NewEmkdir(absPath, err)
		return emkd
	}
	return nil
}

func Clear() error {
	if err := os.RemoveAll(userFilesDir); err != nil {
		emkd := NewEmkdir(userFilesDir, err)
		return emkd
	}

	return nil
}
