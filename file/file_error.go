package file

import (
	"errors"
	"fmt"
)

var ErrFileDoesNotExist = errors.New("file does not exists")
var ErrCannotOpenFile = errors.New("error opening file")
var ErrCannotCreateFile = errors.New("error creating file")
var ErrFileExists = errors.New("file exists")
var ErrTruncated = errors.New("error while truncating file")
var ErrForbbidenAction = errors.New("file forbbiden to remove")
var ErrCouldntWriteFile = errors.New("could not write to file")

type ErrorMkdir struct {
	path string
	err  error
}

func NewEmkdir(abs_path string, e error) ErrorMkdir {
	return ErrorMkdir{
		path: abs_path,
		err:  e,
	}
}

func (em ErrorMkdir) Error() string {
	return fmt.Sprintf("failed to create or remove directory %q: %w", em.path, em.err)
}
