package file

import (
	"os"
	"path/filepath"
)

func HandleTouch(file_name string) error {
	os.Chdir(userFilesDir)
	defer os.Chdir("..") //back to normal(necessay for cleanup)
	var exists bool = doesFileExists(file_name)

	if exists {
		return ErrFileExists
	}

	file, _ := os.Create(file_name)

	defer file.Close()
	return nil
}

func HandleTruncate(file_name string) error {
	var exists bool = doesFileExists(file_name)

	if !exists {
		return ErrFileDoesNotExist
	}

	file, err := os.OpenFile(file_name, os.O_RDWR, 0644)
	if err != nil {
		return ErrCannotOpenFile

	}
	defer file.Close()

	err = file.Truncate(0)
	if err != nil {
		return ErrTruncated
	}

	return nil

}

func HandleRm(file_name string) error {
	var exists bool = doesFileExists(file_name)

	if !exists {
		return ErrFileDoesNotExist
	}

	program_dir, _ := os.Getwd()
	file_path := filepath.Join(program_dir, file_name)
	if isForbidden(file_path) {
		return ErrForbbidenAction
	}
	err := os.Remove(file_path)
	if err != nil {
		return err
	}
	return nil
}
