package reader

import (
	"cli_interpreter/file"
	"errors"
	"fmt"
	"io"
	"regexp"
	stdTime "time"
)

func (r *Reader) Execute(command string) error {

	var err error = nil

	if command == " " || command == "" {
		return nil
	}
	if len(command) > 512 {
		return errors.New("command is longer than 512 characters")
	}

	r.parse_input(command)

	err = r.recognize_command()

	r.clear()

	return err
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

func (r *Reader) recognize_command() error {
	var err error
	command, found := convert_to_enum(r.words[0])
	if !found {
		return errors.New("cannot map command")
	}

	if is_zero_arg_command(command) {
		goto check
	}
	if len(r.words) < 2 {
		r.check_for_more_arguments()
	}
check:
	switch command {
	case echo:
		for i := 1; i < len(r.words); i++ {
			fmt.Print(r.words[i])
			if i < len(r.words)-1 {
				fmt.Print(" ")
			}
		}
		fmt.Print("\n")
	case prompt:
		r.Sign = r.words[1]
	case time:
		fmt.Print("System time: ")
		fmt.Println(stdTime.Now().Clock())
	case date:
		fmt.Print("System date: ")
		fmt.Println(stdTime.Now().Date())
	case touch:
		err = file.Handle_touch(r.words[1])
	case truncate:
		err = file.Handle_truncate(r.words[1])
	case rm:
		err = file.Handle_rm(r.words[1])
	case wc:
		copt, found := convert_command_opt(r.words[1])
		if !found {
			return errors.New("unsupported option type")
		}
		err = r.handle_wc(copt)
	default:
		return errors.New("command unrecognized")
	}

	return err
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
