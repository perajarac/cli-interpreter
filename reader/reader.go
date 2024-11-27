package reader

import (
	"fmt"
	"strings"
)

type Reader struct {
	words []string
}

func (r *Reader) Execute(command string) {
	if command == "" || len(command) > 512 {
		return

	}

	fmt.Printf("Received command %s\n", command)
	r.parse_input(command)

}

func (r *Reader) parse_input(command string) {
	for _, word := range strings.Fields(command) {
		r.words = append(r.words, word)
	}
}
