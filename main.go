package main

import (
	"bufio"
	"cli_interpreter/memory"
	read "cli_interpreter/reader"
	"fmt"
	"log"
	"os"

	"github.com/eiannone/keyboard"
)

var mem *memory.Memory = memory.New()

var reader *read.Reader = &read.Reader{
	Sign:    "$",
	Scanner: bufio.NewReader(os.Stdin),
}

func listen_for_arrow_keys(channel chan string) {
	// Initialize keyboard listener
	err := keyboard.Open()
	if err != nil {
		log.Fatal(err)
	}
	defer keyboard.Close()

	// Listen for arrow key presses
	for {
		_, key, err := keyboard.GetKey()
		if err != nil {
			log.Fatal(err)
		}

		if key == keyboard.KeyEsc {
			break
		}

		switch key {
		case keyboard.KeyArrowUp:
			fmt.Println("ovde sam")
			channel <- mem.Up()
		case keyboard.KeyArrowDown:
			channel <- mem.Down()
		default:
			fmt.Println("ovde sam")
		}

	}
}

func read_coomands(channel chan string) {
	var cmd string
	for {
		fmt.Print(reader.Sign)
		cmd = reader.Read_command()
		channel <- cmd
	}
}

func main() {

	clear := func() {
		mem.Clear()
		reader.Clear()
	}

	command := make(chan string)

	go read_coomands(command)
	go listen_for_arrow_keys(command)

	for cmd := range command {
		err := reader.Execute(cmd)
		if err != nil {
			fmt.Println("Error occurred:", err)
		}
		clear()
	}

}
