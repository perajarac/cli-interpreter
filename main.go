package main

import (
	"cli_interpreter/reader"
	"fmt"
)

func main() {
	reader := &reader.Reader{}
	var input string
	for input != "^V" {
		fmt.Print("$")
		fmt.Scanf("%s", &input)
		reader.Execute(input)
	}

}
