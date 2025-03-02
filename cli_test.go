package main

import (
	"strings"
	"testing"

	f "github.com/perajarac/cli-interpreter/file"
	r "github.com/perajarac/cli-interpreter/reader"
)

var reader_test *r.Reader = r.NewReader()

const helpText string = `Available commands:
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

func TestValidateCLI(t *testing.T) {
	r.SetUpUser()
	defer reader_test.Clear()
	tests := []struct {
		command  string
		wantErr  error
		expected string
	}{
		{"echo PeraJarac", f.ErrFileDoesNotExist, ""},
		{"echo Pera Jarac", f.ErrFileDoesNotExist, ""},
		{"echo \"PeraJarac\"", nil, "PeraJarac"},
		{"echo \"Pera Jarac\"", nil, "Pera Jarac"},
		{"prompt %", nil, ""},
		{"time", nil, r.TimeOrDate(3)},
		{"date", nil, r.TimeOrDate(4)},
		{"touch pera.txt", nil, ""},
		{"truncate pera.txt", nil, ""},
		{"rm pera.txt", nil, ""},
		{"wc -c \"Pera\"", nil, "4"},
		{"wc -w \"Pera\"", nil, "1"},
		{"wc -w \"Pera jarac\"", nil, "2"},
		{"tr \"Bleja je lepa\" \"lepa\" \"bleja\"", nil, "Bleja je bleja"},
		{"help", nil, helpText},
		{"version", nil, r.Ver},
		{"version | tr \".\"", nil, "103"},
		{"version | echo", nil, "1.0.3"},
		{"echo \"hello\" | wc -w | tr \"1\" \"one\"", nil, "one"},
		{"echo \"hello world\" | wc -c | tr \"11\" \"eleven\"", nil, "eleven"},
		{"echo \"sample text\" | tr \"sample\" \"example\" | wc -w", nil, "2"},
		{"echo \"test\" | tr \"t\" \"T\" | wc -c", nil, "4"},
		{"echo \"pipe test\" | wc -w | wc -c", nil, "1"},
	}
	for _, tt := range tests {
		ret, err := reader_test.RunCommand(tt.command)
		ret = strings.TrimSuffix(ret, "\n")
		if ((err != nil) && err != tt.wantErr) || (ret != tt.expected) {
			t.Errorf("Error occured: command = %v output = %v, expected = %v, error = %v, wantErr = %v", tt.command, ret, tt.expected, err, tt.wantErr)
		}
	}
}
