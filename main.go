package main

import (
	read "github.com/perajarac/cli_interpreter/reader"
)

var reader *read.Reader = read.NewReader()

func main() {
	for {
		reader.MainLoop()
	}
}
