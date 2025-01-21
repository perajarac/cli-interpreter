package main

import (
	read "cli_interpreter/reader"
)

var reader *read.Reader = read.NewReader()

func main() {
	for {
		reader.MainLoop()
	}
}
