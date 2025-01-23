package reader

import (
	"bufio"
	"cli_interpreter/memory"
	"fmt"
	"io"
	"os"
	"strings"
)

type Reader struct {
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

func (r *Reader) MainLoop() {
	var ret string
	var err error
	fmt.Print(r.Sign)
	command := r.ReadCommand()

	ret, err = r.RunCommand(command)

	if err != nil {
		fmt.Println("Error occured: ", err)
	}
	if ret != "" {
		fmt.Println(ret)
	}
	r.Clear()
}

func (r *Reader) RunCommand(cmd string) (string, error) {
	cmd = strings.TrimSpace(cmd)
	if cmd == "" {
		return "", nil
	}
	if len(cmd) > 512 {
		return "", ErrToLongCommand
	}

	if strings.Contains(cmd, "|") {
		return r.handlePipes(cmd)
	} else {
		return r.handleSimpleCmd(cmd)
	}

}

func (r *Reader) ReadCommand() string {
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
