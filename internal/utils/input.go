package utils

import (
	"bufio"
	"fmt"
	"os"
)

func GetInput() string {
	fmt.Print("> ")
	scanner := bufio.NewScanner(os.Stdin)
	scanned := scanner.Scan()
	if !scanned {
		return ""
	}
	return scanner.Text()
}
