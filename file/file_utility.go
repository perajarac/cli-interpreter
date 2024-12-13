package file

import (
	"os"
	"path/filepath"
	"strings"
)

var forbbiden_extensions []string = []string{".go", ".exe", ".dll", ".sh"}

func does_file_exists(file_name string) bool {
	_, error := os.Stat(file_name)
	return !os.IsNotExist(error)

}

func is_forbbiden(file_path string) bool {
	extension := strings.ToLower(filepath.Ext(file_path))
	for _, ext := range forbbiden_extensions {
		if extension == ext {
			return true
		}
	}
	return false
}
