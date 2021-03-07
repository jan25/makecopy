package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

func prompt(q string, suggest string) (string, error) {
	p := ""
	if suggest != "" {
		p = fmt.Sprintf("%s (%s): ", q, suggest)
	} else {
		p = fmt.Sprintf("%s: ", q)
	}

	reader := bufio.NewReader(os.Stdin)
	fmt.Print(p)
	result, err := reader.ReadString('\n')
	return strings.TrimSpace(result), err
}
