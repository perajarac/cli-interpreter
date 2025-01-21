package reader

import (
	"cli_interpreter/file"
	"fmt"
	"io"
	"strings"
	stdTime "time"
)

func (r *Reader) Execute(command string) error {

	var err error = nil

	command = strings.TrimSpace(command) // if

	if command == "" {
		return nil
	}
	if len(command) > 512 {
		return ErrToLongCommand
	}

	r.parse_input(command)
	err = r.recognize_command()

	return err
}

func (r *Reader) recognize_command() error {
	var err error
	command, found := convert_to_enum(r.words[0])
	if !found {
		return ErrCannotMapCommand
	}

	if is_zero_arg_command(command) {
		goto check
	}
	if len(r.words) < 2 {
		if command == wc || command == tr || command == head {
			return ErrInvalidFormat
		}
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
			return ErrUnsupportedOptionType
		}
		ret := r.handle_wc(copt)
	case tr:
		ret, err = r.handle_tr()
	case batch: //TODO: make this work
		for _, v := range r.words {
			fmt.Println(v)
		}
	case help:
		Help()
	case version:
		Version()

	default:
		return ErrCannotMapCommand
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

func (r *Reader) MainLoop() {
	fmt.Print(r.Sign)
	command := r.Read_command()
	err := r.Execute(command)
	if err != nil {
		fmt.Println("Error occured: ", err)
	}
	r.Clear()
}
