package reader

import (
	"bufio"
	"cli_interpreter/memory"
	"fmt"
	"os"
	"regexp"
	"strings"
	stdTime "time"
)

const ver = "1.0.0"

type Reader struct {
	Sign    string
	Scanner *bufio.Reader
	Memmory *memory.Memory
}

type Command struct {
	words []string
	ct    command_type
	opt   command_option
	arg   string
}

func NewReader() *Reader {
	return &Reader{
		Sign:    "$",
		Scanner: bufio.NewReader(os.Stdin),
		Memmory: memory.New(),
	}
}

func NewCommand() *Command {
	return &Command{
		words: []string{},
		arg:   "",
	}
}

type command_type int
type command_option int

const (
	unkct command_option = iota
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

// TODO: make this work
// func (r *Reader) batch_helper(command string) []string {
// 	words := strings.Fields(command)

// 	if len(words) == 1 && words[0] == "batch" {
// 		return []string{"batch"}
// 	}

// 	allCommands := strings.Join(words[1:], " ")
// 	parts := strings.Split(allCommands, ";")

// 	for i := range parts {
// 		parts[i] = strings.TrimSpace(parts[i])
// 	}
// 	if words[0] != "batch" && r.words[0] == "batch" {
// 		return parts
// 	} else if words[0] == "batch" {
// 		return append([]string{"batch"}, parts...)
// 	}

// 	return nil
// }

func Version() string {
	return ver
}

func Help() string {
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

	return helpText
}

// function that fills all fields of c *Command
func (c *Command) parse_input(command string) error {
	re := regexp.MustCompile(`-[A-Z][a-z]*|\[[^\]]*\]|"[^"]*"|\S+`)
	var found bool
	matches := re.FindAllString(command, -1)
	for _, word := range matches {
		if word[0] == '"' || word[0] == '[' {
			word = word[1 : len(word)-1]
		}
		c.words = append(c.words, word)
	}
	c.ct, found = getCommandType(c.words[0])
	if !found {
		return ErrCannotMapCommand
	}
	if len(c.words) > 1 {
		if len(c.words[1]) > 0 && c.words[1][0] == '-' {
			c.opt, found = getCommandOpt(c.words[1])
			if !found {
				return ErrUnsupportedOptionType
			}
		}
		if c.opt == unkct { //for commands with no optional flag
			c.arg = strings.Join(c.words[1:], " ")
		} else {
			c.arg = strings.Join(c.words[2:], " ")
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

func (r *Reader) check_for_more_arguments(c *Command) {
	if c.words[0] == "tr" {
		return
	}
	more_args := r.Read_command()
	c.parse_input(more_args)
}

func (r *Reader) Clear() { //TODO: add more stuff if necessary
	r.Memmory.Clear()
}

func is_zero_arg_command(command command_type) bool {
	return command == time || command == date || command == version || command == help
}

func count_letters(word string) int {
	return len(word)
}

func count_words(sentence string) int {
	words := strings.Fields(sentence)
	return len(words)
}

func (r *Reader) handle_wc(copt command_option, comm *Command) int {
	if len(comm.words) < 3 {
		r.check_for_more_arguments(comm)
	}
	var ret int = 0

	for i := 2; i < len(comm.words); i++ {
		if i >= 3 && copt == c {
			ret += 1
		}
		if copt == w {
			ret += count_words(comm.words[i])
		} else {
			ret += count_letters(comm.words[i])
		}
	}

	return ret

}

func (r *Reader) handle_tr(c *Command) (string, error) {
	var ret string = " "
	if len(c.words) < 3 {
		return ret, ErrToFewArgs
	}

	if len(c.words) > 4 {
		return ret, ErrInvalidFormat
	}

	if len(c.words) == 3 {
		ret = strings.ReplaceAll(c.words[1], c.words[2], "")
		return ret, nil
	}
	ret = strings.ReplaceAll(c.words[1], c.words[2], c.words[3])

	return ret, nil
}

func Echo(c *Command) string {
	var ret string
	for i := 1; i < len(c.words); i++ {
		ret = ret + c.words[i]
		if i < len(c.words)-1 {
			ret = ret + " "
		}
	}
	ret += "\n"

	return ret
}

func TimeOrDate(ct command_type) string {
	var ret string
	var timeString string
	if ct == time {
		ret = ret + "System time: "
		first, second, third := stdTime.Now().Clock()
		timeString = fmt.Sprintf("%02d:%02d:%02d", first, second, third)
	} else {
		ret = ret + "System date: "
		first, second, third := stdTime.Now().Date()
		timeString = fmt.Sprintf("%02d:%02d:%02d", first, second, third)
	}
	ret = ret + timeString + "\n"
	return ret
}
