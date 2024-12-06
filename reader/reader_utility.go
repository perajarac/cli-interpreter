package reader

import (
	"strings"
)

type command_type int

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
)

var CommandMap = map[string]command_type{
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
}

func ConvertToEnum(word string) (command_type, bool) {
	word = strings.ToLower(word)
	cmd, found := CommandMap[word]
	return cmd, found
}

type Reader struct {
	words []string
}
