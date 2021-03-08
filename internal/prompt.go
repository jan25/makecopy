package internal

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

// Prompt prompts a question to user.
func Prompt(q string, suggest string) (string, error) {
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
