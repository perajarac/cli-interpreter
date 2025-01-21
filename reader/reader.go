package reader

import (
	"cli_interpreter/file"
	"fmt"
	"io"
	"strings"
)

func (r *Reader) Execute(command string) (string, error) {

	var err error = nil
	var ret string = ""

	command = strings.TrimSpace(command) // if

	if command == "" {
		return ret, nil
	}
	if len(command) > 512 {
		return ret, ErrToLongCommand
	}

	var comm *Command = NewCommand()
	err = comm.parse_input(command)
	if err != nil {
		return ret, err
	}
	ret, err = r.recognize_command(comm)

	return ret, err
}

func (r *Reader) recognize_command(comm *Command) (string, error) {
	var err error = nil
	var ret string = ""

	if is_zero_arg_command(comm.ct) {
		goto check
	}
	if len(comm.words) < 2 {
		if comm.ct == wc || comm.ct == tr || comm.ct == head {
			return ret, ErrInvalidFormat
		}
		r.check_for_more_arguments(comm)
	}
check:
	switch comm.ct {
	case echo:
		ret = Echo(comm)
	case prompt:
		r.Sign = comm.words[1]
	case time:
		ret = TimeOrDate(time)
	case date:
		ret = TimeOrDate(date)
	case touch:
		err = file.Handle_touch(comm.arg)
	case truncate:
		err = file.Handle_truncate(comm.arg)
	case rm:
		err = file.Handle_rm(comm.arg)
	case wc:
		ret = fmt.Sprintf("%d", r.handle_wc(comm.opt, comm))
	case tr:
		ret, err = r.handle_tr(comm)
	case batch: //TODO: make this work
		for _, v := range comm.words {
			fmt.Println(v)
		}
	case help:
		ret = Help()
	case version:
		ret = Version()

	default:
		return ret, ErrCannotMapCommand
	}

	return ret, err
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
	ret, err := r.Execute(command)
	if err != nil {
		fmt.Println("Error occured: ", err)
	}
	fmt.Print(ret)
	r.Clear()
}
