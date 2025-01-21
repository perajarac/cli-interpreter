package reader

import (
	"bufio"
	"cli_interpreter/memory"
	"fmt"
	"os"
	"regexp"
	"strings"
)

const ver = "1.0.0"

type Reader struct {
	words   []string
	Sign    string
	Scanner *bufio.Reader
	Memmory *memory.Memory
}

func NewReader() *Reader {
	return &Reader{
		Sign:    "$",
		Scanner: bufio.NewReader(os.Stdin),
		Memmory: memory.New(),
	}
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

func (r *Reader) batch_helper(command string) []string {
	words := strings.Fields(command)

	if len(words) == 1 && words[0] == "batch" {
		return []string{"batch"}
	}

	allCommands := strings.Join(words[1:], " ")
	parts := strings.Split(allCommands, ";")

	for i := range parts {
		parts[i] = strings.TrimSpace(parts[i])
	}
	if words[0] != "batch" && r.words[0] == "batch" {
		return parts
	} else if words[0] == "batch" {
		return append([]string{"batch"}, parts...)
	}

	return nil
}

func Version() {
	fmt.Println(ver)
}

func Help() {
	fmt.Println("Available commands:")
	fmt.Println("1. echo [argument]")
	fmt.Println("   - Sends the input string directly to the output without any modifications.")

	fmt.Println("2. prompt [argument]")
	fmt.Println("   - Sets the command prompt to the specified string argument.")

	fmt.Println("3. time")
	fmt.Println("   - Outputs the current system time.")

	fmt.Println("4. date")
	fmt.Println("   - Outputs the current system date.")

	fmt.Println("5. touch [filename]")
	fmt.Println("   - Creates an empty file with the specified filename in the current directory.")
	fmt.Println("     Outputs an error message if the file already exists.")

	fmt.Println("6. truncate [filename]")
	fmt.Println("   - Deletes the content of the specified file in the current directory.")

	fmt.Println("7. rm [filename]")
	fmt.Println("   - Removes the specified file from the file system in the current directory.")

	fmt.Println("8. wc -opt [argument]")
	fmt.Println("   - Counts words or characters in the input text based on the option.")
	fmt.Println("     -w for words, -c for characters.")

	fmt.Println("9. tr [argument] what [with]")
	fmt.Println("   - Replaces all occurrences of the string 'what' with the string 'with' in the input text.")
	fmt.Println("     If 'with' is not specified, 'what' will be removed.")

	fmt.Println("10. head -ncount [argument]")
	fmt.Println("    - Outputs the first 'count' lines of the input text.")

	fmt.Println("11. batch [argument]")
	fmt.Println("    - Interprets multiple command lines from the input as if they were entered one by one in the terminal.")

	fmt.Println("12. help")
	fmt.Println("    - Displays the documentation for all available commands.")
	fmt.Println("13. version")
	fmt.Println("    - Displays the version of the program.")
}

func (r *Reader) parse_input(command string) {
	// temp := r.batch_helper(command) #TODO: change this

	// r.words = append(r.words, temp...)

	// if r.words != nil {
	// 	return
	// }

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
	if r.words[0] == "tr" {
		return
	}
	more_args := r.Read_command()
	r.parse_input(more_args)
}

func (r *Reader) Clear() { //TODO: add more stuff if necessary
	if len(r.words) > 0 {
		r.words = r.words[:0]
	}
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

func (r *Reader) handle_wc(copt command_option) int {
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

	return ret

}

func (r *Reader) handle_tr() (string, error) {
	var ret string = " "
	if len(r.words) < 3 {
		return ret, ErrToFewArgs
	}

	if len(r.words) > 4 {
		return ret, ErrInvalidFormat
	}

	if len(r.words) == 3 {
		ret = strings.ReplaceAll(r.words[1], r.words[2], "")
		return ret, nil
	}
	ret = strings.ReplaceAll(r.words[1], r.words[2], r.words[3])

	return ret, nil
}
