package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

func Retrieve_user_input() (string, error) {

	reader := bufio.NewReader(os.Stdin)

	target, err := reader.ReadString('\n')
	if err != nil {
		msg := "Error reading input:"
		fmt.Println(msg, err)
		return target, err
	}
	return strings.TrimSpace(target), nil
}
