package reader

import (
	"errors"
	"fmt"
	"regexp"
)

func (r *Reader) Execute(command string) error {
	if command == " " || command == "" {
		return nil
	}

	if len(command) > 512 {
		return errors.New("Command is longer than 512 characters!")
	}

	r.parse_input(command)

	if !recognize_command(r) {
		return errors.New("Command unrecognized")
	}

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

func recognize_command(r *Reader) bool {
	command, found := ConvertToEnum(r.words[0])
	if !found {
		return false
	}

	switch command {
	case echo:
		fmt.Println(r.words[1])
	default:
		fmt.Println(command)
	}

	return true
}
