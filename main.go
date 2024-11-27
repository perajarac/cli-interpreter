package main

import (
	"bufio"
	"cli_interpreter/reader"
	"fmt"
	"os"
)

func main() {
	reader := &reader.Reader{}
	scanner := bufio.NewReader(os.Stdin)
	for {
		fmt.Print("$")
		input, err := scanner.ReadString('\n')
		if err != nil {
			fmt.Println("Error while reading", err)
			return
		}
		reader.Execute(input[:len(input)-1])
	}
}
