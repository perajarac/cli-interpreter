package reader

import (
	"errors"
	"fmt"
	"io"
	"regexp"
)

func (r *Reader) Execute(command string) error {
	if command == " " || command == "" {
		return nil
	}
	if len(command) > 512 {
		return errors.New("command is longer than 512 characters")
	}

	r.parse_input(command)

	// fmt.Println("\t Current state of system: \t", r.words, r.Sign)

	if !r.recognize_command() {
		r.clear()
		return errors.New("command unrecognized") //TODO: make generic error
	}

	r.clear()

	return nil
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

func (r *Reader) recognize_command() bool {
	command, found := convertToEnum(r.words[0])
	if !found {
		return false
	}
	if len(r.words) < 2 {
		r.check_for_more_arguments()
	}
	switch command {
	case echo:
		fmt.Println(r.words[1])
	case prompt:
		r.Sign = r.words[1]
	default:
	}

	return true
}

func (r *Reader) Read_command() string {
	input, err := r.Scanner.ReadString('\n')
	if err != nil {
		if err == io.EOF {
			return ""
		}
		fmt.Println("Error while reading", err)
		return ""
	}

	return input[:len(input)-1]
}
