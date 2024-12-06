package reader

import (
	"bufio"
	"strings"
)

type command_type int

const (
	echo command_type = iota + 1
	prompt
	time
	date
	touch
	truncate
	rm
	wc
	tr
	head
	batch
)

var CommandMap = map[string]command_type{
	"echo":     echo,
	"prompt":   prompt,
	"time":     time,
	"date":     date,
	"touch":    touch,
	"truncate": truncate,
	"rm":       rm,
	"wc":       wc,
	"tr":       tr,
	"head":     head,
	"batch":    batch,
}

func convertToEnum(word string) (command_type, bool) {
	word = strings.ToLower(word)
	cmd, found := CommandMap[word]
	return cmd, found
}

type Reader struct {
	words   []string
	Sign    string
	Scanner *bufio.Reader
}

func (r *Reader) check_for_more_arguments() {
	if len(r.words) == 1 {
		more_args := r.Read_command()
		r.parse_input(more_args)
		return
	}

	if len(r.words) == 2 && r.words[1][0] == '-' {
		return //TODO: add checking for arguments instructions and make it more robust
	}
}

func (r *Reader) clear() { //TODO: add more stuff if necessary
	if len(r.words) > 0 {
		r.words = r.words[:0]
	}
}
