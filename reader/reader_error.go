package reader

import "errors"

var ErrToLongCommand = errors.New("command is longer than 512 characters")
var ErrCannotMapCommand = errors.New("cannot map command")
var ErrInvalidFormat = errors.New("invalid command format")
var ErrUnsupportedOptionType = errors.New("unsupported option type")
var ErrToFewArgs = errors.New("to few arguments")
var ErrNoInputFile = errors.New("no input file provided")
