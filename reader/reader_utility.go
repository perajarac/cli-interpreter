package reader

import (
	"fmt"
	"regexp"
	"strings"
	stdTime "time"
	"unicode/utf8"

	"github.com/perajarac/cli-interpreter/file"
)

const Ver = "1.0.2"

type Command struct {
	words  []string
	ct     command_type
	opt    command_option
	arg    string
	output string
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
}

var command_opt_map = map[string]command_option{
	"-w": w,
	"-c": c,
	"-n": n,
}

func Version(sentence string) string {
	if sentence != "" {
		return sentence + " " + Ver
	}
	return Ver
}

func Help(sentence string) string {
	helpText := `Available commands:
	1. echo [argument]
	   - Sends the input string directly to the output without any modifications.
	
	2. prompt [argument]
	   - Sets the command prompt to the specified string argument.
	
	3. time
	   - Outputs the current system time.
	
	4. date
	   - Outputs the current system date.
	
	5. touch [filename]
	   - Creates an empty file with the specified filename in the current directory.
		 Outputs an error message if the file already exists.
	
	6. truncate [filename]
	   - Deletes the content of the specified file in the current directory.
	
	7. rm [filename]
	   - Removes the specified file from the file system in the current directory.
	
	8. wc -opt [argument]
	   - Counts words or characters in the input text based on the option.
		 -w for words, -c for characters.
	
	9. tr [argument] what [with]
	   - Replaces all occurrences of the string 'what' with the string 'with' in the input text.
		 If 'with' is not specified, 'what' will be removed.
	
	10. head -ncount [argument]
		- Outputs the first 'count' lines of the input text.
	
	11. batch [argument]
		- Interprets multiple command lines from the input as if they were entered one by one in the terminal.
	
	12. help
		- Displays the documentation for all available commands.
	
	13. version
		- Displays the version of the program.`

	if sentence != "" {
		return sentence + " " + helpText
	}
	return helpText
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
	var file_content string
	c.words, file_content, err = file.CheckArgument(c.words)
	if err != nil {
		return err
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

func (r *Reader) HandleWc(copt command_option, comm *Command) (int, error) {
	if comm.opt == unkco {
		return 0, ErrInvalidFormat
	}
	comm.arg = strings.ReplaceAll(comm.arg, `"`, "")
	if comm.arg == "" {
		r.checkForMoreArgs(comm)
	}
	var ret int = 0

	if copt == w {
		ret = countWords(comm.arg)
	} else {
		ret = countLetters(comm.arg)
	}

	return ret, nil

}

func (r *Reader) HandleTr(c *Command) (string, error) {
	reg := regexp.MustCompile(`"([^"]*)"`)
	matches := reg.FindAllString(c.arg, -1)

	if len(matches) < 2 || len(matches) > 3 {
		return "", ErrInvalidFormat
	}
	var ret string = ""

	if len(matches) == 2 {
		ret = strings.ReplaceAll(matches[0], strings.Trim(matches[1], `"`), "")
	} else {
		ret = strings.ReplaceAll(matches[0], strings.Trim(matches[1], `"`), strings.Trim(matches[2], `"`))
	}

	return ret, nil
}

func Echo(c *Command) string {
	return strings.ReplaceAll(c.arg, `"`, "")
}

func TimeOrDate(ct command_type) string {
	var ret string
	var timeString string
	if ct == time {
		first, second, third := stdTime.Now().Clock()
		timeString = fmt.Sprintf("%02d:%02d:%02d", first, second, third)
	} else {
		first, second, third := stdTime.Now().Date()
		timeString = fmt.Sprintf("%02d:%02d:%02d", first, second, third)
	}
	ret = ret + timeString
	return ret
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
