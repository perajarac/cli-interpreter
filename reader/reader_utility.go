package reader

import (
	"fmt"
	"regexp"
	"strings"
	"unicode/utf8"

	"github.com/perajarac/cli-interpreter/file"
)

const Ver = "1.0.3"

type Command struct {
	words  []string
	ct     command_type
	opt    command_option
	arg    string
	output string
	input  bool
}

func NewCommand() *Command {
	return &Command{
		words:  []string{},
		arg:    "",
		output: "",
	}
}

type command_type int
type command_option int

const (
	unkco command_option = iota
	w
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
	help
	version
	cat
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
	"help":     help,
	"version":  version,
	"cat":      cat,
}

var command_opt_map = map[string]command_option{
	"-w": w,
	"-c": c,
	"-n": n,
}

// function that fills all fields of c *Command
func (c *Command) parseInput(command string) error {
	re := regexp.MustCompile(`-[A-Z][a-z]*|"[^"]*"|\S+`)
	var found bool
	matches := re.FindAllString(command, -1)

	c.words = append(c.words, matches...)
	c.ct, found = getCommandType(c.words[0])
	if !found {
		return ErrCannotMapCommand
	}

	var err error

	c.checkOutputPath()

	//check if command has file(< or basic file.txt) as a argument if has, argument becomes all file content
	var file_content string = ""
	c.words, file_content, err = file.CheckArgument(c.words)
	if err != nil {
		return err
	}
	if file_content != "" {
		c.input = true
	}

	c.arg = c.arg + file_content

	if len(c.words) > 1 {
		//if command has flag
		if len(c.words[1]) > 0 && c.words[1][0] == '-' {
			c.opt, found = getCommandOpt(c.words[1])
			if !found {
				return ErrUnsupportedOptionType
			}
		}
		if c.opt == unkco { //for commands with no optional flag
			c.arg = c.arg + strings.Join(c.words[1:], " ")
		} else {
			c.arg = c.arg + strings.Join(c.words[2:], " ")
		}
	}

	return nil
}

func getCommandType(word string) (command_type, bool) {
	word = strings.ToLower(word)
	cmd, found := command_map[word]
	return cmd, found
}

func getCommandOpt(word string) (command_option, bool) {
	word = strings.ToLower(word)
	cmd, found := command_opt_map[word]
	return cmd, found
}

func (r *Reader) checkForMoreArgs(c *Command) {
	if c.ct == tr {
		return
	}
	more_args := r.ReadCommand()
	c.parseInput(more_args)
}

func (r *Reader) Clear() { //TODO: add more stuff if necessary
	err := file.Clear()
	if err != nil {
		panic(err)
	}
	r.Memmory.Clear()
}

func isZeroArg(command command_type) bool {
	return command == time || command == date || command == version || command == help || command == echo
}

func countLetters(sentence string) int {
	return utf8.RuneCountInString(sentence)
}

func countWords(sentence string) int {
	words := strings.Fields(sentence)
	return len(words)
}

func (r *Reader) handlePipes(cmd string) (string, error) {
	commands := strings.Split(cmd, "|")
	var last string
	var err error

	for _, comm := range commands {
		command := strings.TrimSpace(comm)
		if command == "" {
			return "", nil
		}
		if len(command) > 512 {
			return "", ErrToLongCommand
		}

		commObj := NewCommand()
		commObj.arg = last
		err = commObj.parseInput(command)
		if err != nil {
			return "", err
		}
		last, err = r.recognizeCommand(commObj)
		if err != nil {
			return "", err
		}
		if commObj.output != "" {
			err = file.WriteOutput(commObj.output, last)
			if err != nil {
				return "", err
			}
			last = ""
			continue
		}
		last = "\"" + last + "\""
	}
	return strings.Trim(last, `"`), nil
}

func (r *Reader) handleSimpleCmd(cmd string) (string, error) {
	commObj := NewCommand()
	err := commObj.parseInput(cmd)
	if err != nil {
		return "", err
	}
	ret, err := r.recognizeCommand(commObj)
	if err != nil {
		return "", err
	}
	ret = strings.Trim(ret, `"`)

	if commObj.output != "" {
		err = file.WriteOutput(commObj.output, ret)
		if err != nil {
			return "", err
		}
		ret = ""
	}

	return ret, nil
}

func (r *Reader) recognizeCommand(comm *Command) (string, error) {
	var err error = nil
	var ret string = ""

	if comm.arg == "" && !isZeroArg(comm.ct) {
		if comm.ct == wc || comm.ct == tr || comm.ct == head {
			return ret, ErrInvalidFormat
		}
		r.checkForMoreArgs(comm)
	}
	switch comm.ct {
	case echo:
		ret = Echo(comm)
	case prompt:
		r.Sign = comm.arg
	case time:
		ret = TimeOrDate(time)
	case date:
		ret = TimeOrDate(date)
	case touch:
		err = file.HandleTouch(comm.arg)
	case truncate:
		err = file.HandleTruncate(comm.arg)
	case rm:
		err = file.HandleRm(comm.arg)
	case wc:
		var temp int
		temp, err = r.HandleWc(comm.opt, comm)
		ret = fmt.Sprintf("%d", temp)
	case tr:
		ret, err = r.HandleTr(comm)
	case batch: //TODO:make this work
		for _, v := range comm.words {
			fmt.Println(v)
		}
	case help:
		ret = Help(comm.arg)
	case version:
		ret = Version(comm.arg)
	case cat:
		ret, err = Cat(comm)

	default:
		return ret, ErrCannotMapCommand
	}

	return ret, err
}

func (cmd *Command) checkOutputPath() {
	index := len(cmd.words) - 1
	if index > 0 {
		if cmd.words[index][0] == '>' {
			cmd.output = strings.ReplaceAll(cmd.words[index], ">", "")
			cmd.words = file.RemoveAtIndex(cmd.words, index)
		}
	}
}

func SetUpUser() {
	if err := file.EnsureUserFilesDir(); err != nil {
		panic(err)
	}
}

// func CleanUser(){

// }
