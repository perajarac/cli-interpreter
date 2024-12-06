package main

import (
	"bufio"
	"cli_interpreter/reader"
	"fmt"
	"os"
)

func main() {
	reader := &reader.Reader{
		Sign:    "$",
		Scanner: bufio.NewReader(os.Stdin),
	}

	for {
		fmt.Print(reader.Sign)
		command := reader.Read_command()
		err := reader.Execute(command)
		if err != nil {
			fmt.Println("Error occured", err)
		}
	}
}
