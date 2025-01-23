package file

import (
	"errors"
	"os"
	"path/filepath"
)

func HandleTouch(file_name string) error {
	var exists bool = does_file_exists(file_name)

	if exists {
		return errors.New("file exists")
	}

	file, _ := os.Create(file_name)

	defer file.Close()
	return nil
}

func HandleTruncate(file_name string) error {
	var exists bool = does_file_exists(file_name)

	if !exists {
		return errors.New("file does not exists")
	}

	file, err := os.OpenFile(file_name, os.O_RDWR, 0644)
	if err != nil {
		return errors.New("error opening file")

	}
	defer file.Close()

	err = file.Truncate(0)
	if err != nil {
		return errors.New("error while truncating file")
	}

	return nil

}

func HandleRm(file_name string) error {
	var exists bool = does_file_exists(file_name)

	if !exists {
		return errors.New("file does not exists")
	}

	program_dir, _ := os.Getwd()
	file_path := filepath.Join(program_dir, file_name)
	if is_forbbiden(file_path) {
		return errors.New("file forbbiden to remove")
	}
	err := os.Remove(file_path)
	if err != nil {
		return err
	}
	return nil
}
