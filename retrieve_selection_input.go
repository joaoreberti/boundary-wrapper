package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
)

func Retrieve_selection_input() (int, error) {
	reader := bufio.NewReader(os.Stdin)
	option, err := reader.ReadString('\n')
	if err != nil {
		msg := "Error reading input:"
		fmt.Println(msg, err)
		return 0, nil
	}
	number, err := strconv.ParseUint(strings.TrimSpace(option), 10, 32)
	finalIntNum := int(number) //Convert uint64 To int
	return finalIntNum, err
}
