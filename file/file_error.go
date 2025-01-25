package file

import "errors"

var ErrFileDoesNotExist = errors.New("file does not exists")
var ErrCannotOpenFile = errors.New("error opening file")
var ErrCannotCreateFile = errors.New("error creating file")
var ErrFileExists = errors.New("file exists")
var ErrTruncated = errors.New("error while truncating file")
var ErrForbbidenAction = errors.New("file forbbiden to remove")
var ErrCouldntWriteFile = errors.New("could not write to file")
