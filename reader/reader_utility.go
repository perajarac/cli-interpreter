package reader

import (
	"bufio"
	"fmt"
	"regexp"
	"strings"
)

type Reader struct {
	words   []string
	Sign    string
	Scanner *bufio.Reader
}

type command_type int
type command_option int

const (
	w command_option = iota + 1
	c
	n
)

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

var command_map = map[string]command_type{
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

var command_opt_map = map[string]command_option{
	"-w": w,
	"-c": c,
	"-n": n,
}

func (r *Reader) parse_input(command string) {
	re := regexp.MustCompile(`-[A-Z][a-z]*|\[[^\]]*\]|"[^"]*"|\S+`)
	matches := re.FindAllString(command, -1)
	for _, word := range matches {
		if word[0] == '"' || word[0] == '[' {
			word = word[1 : len(word)-1]
		}
		r.words = append(r.words, word)
	}
}

func convert_to_enum(word string) (command_type, bool) {
	word = strings.ToLower(word)
	cmd, found := command_map[word]
	return cmd, found
}

func convert_command_opt(word string) (command_option, bool) {
	word = strings.ToLower(word)
	cmd, found := command_opt_map[word]
	return cmd, found
}

func (r *Reader) check_for_more_arguments() {
	more_args := r.Read_command()
	r.parse_input(more_args)
}

func (r *Reader) Clear() { //TODO: add more stuff if necessary
	if len(r.words) > 0 {
		r.words = r.words[:0]
	}
}

func is_zero_arg_command(command command_type) bool {
	return command == time || command == date
}

func count_letters(word string) int {
	return len(word)
}

func count_words(sentence string) int {
	words := strings.Fields(sentence)
	return len(words)
}

func (r *Reader) handle_wc(copt command_option) error {
	if len(r.words) < 3 {
		r.check_for_more_arguments()
	}
	var ret int = 0

	for i := 2; i < len(r.words); i++ {
		if i >= 3 && copt == c {
			ret += 1
		}
		if copt == w {
			ret += count_words(r.words[i])
		} else {
			ret += count_letters(r.words[i])
		}
	}

	if copt == w {
		fmt.Println("Number of words: ", ret)
	} else {
		fmt.Println("Number of letters: ", ret)
	}

	return nil

}
