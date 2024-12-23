package main

import (
	"bufio"
	"cli_interpreter/memory"
	read "cli_interpreter/reader"
	"fmt"
	"os"
)

var mem *memory.Memory = memory.New()

var reader *read.Reader = &read.Reader{
	Sign:    "$",
	Scanner: bufio.NewReader(os.Stdin),
}

func main() {

	clear := func() {
		mem.Clear()
		reader.Clear()
	}

	for {
		fmt.Print(reader.Sign)
		command := reader.Read_command()
		err := reader.Execute(command)
		if err != nil {
			fmt.Println("Error occured: ", err)
		}
		clear()
		mem.Push(command)
	}

}
