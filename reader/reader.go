package reader

import "fmt"

type Reader struct {
	Input string
}

func (r *Reader) Execute(command string) {
	r.Input = command
	if command == "" || len(command) > 512 {
		return

	}
	fmt.Printf("Received command:%s\n", command)
}
