package file

import (
	"os"
	"path/filepath"
	"strings"
)

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
	file, err := os.OpenFile(outputFile, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0644)
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
	for i, word := range words[1:] {

		if strings.Contains(word, "\"") || word[0] == '>' || word[0] == '-' {
			continue
		}
		words = RemoveAtIndex(words, i)
		word = strings.ReplaceAll(word, "<", "")
		arg, err := readFromFile(word)
		if err != nil {
			return words, "", ErrFileDoesNotExist
		}
		return words, arg, nil
	}
	return words, "", nil
}

func RemoveAtIndex(s []string, i int) []string {
	if i < 0 || i >= len(s) {
		return s
	}
	return append(s[:i], s[i+1:]...)
}
