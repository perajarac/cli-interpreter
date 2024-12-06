package main

import (
	"bufio"
	"cli_interpreter/reader"
	"fmt"
	"io"
	"os"
)

func main() {
	for {
		reader := &reader.Reader{}
		scanner := bufio.NewReader(os.Stdin)
		fmt.Print("$")
		input, err := scanner.ReadString('\n')
		if err != nil {
			if err == io.EOF {
				return
			}
			fmt.Println("Error while reading", err)
			return
		}
		reader.Execute(input[:len(input)-1])
	}
}
