package reader

import (
	"fmt"
	"regexp"
	"strings"
	stdTime "time"
)

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

func Cat(comm *Command) (string, error) {
	if !comm.input {
		return "", ErrNoInputFile
	}
	return Echo(comm), nil
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

	5. cat
		- Shows file content on terminal or write it in given file
	
	6. touch [filename]
	   - Creates an empty file with the specified filename in the current directory.
		 Outputs an error message if the file already exists.
	
	7. truncate [filename]
	   - Deletes the content of the specified file in the current directory.
	
	8. rm [filename]
	   - Removes the specified file from the file system in the current directory.
	
	9. wc -opt [argument]
	   - Counts words or characters in the input text based on the option.
		 -w for words, -c for characters.
	
	10. tr [argument] what [with]
	   - Replaces all occurrences of the string 'what' with the string 'with' in the input text.
		 If 'with' is not specified, 'what' will be removed.
	
	
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
