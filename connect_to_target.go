package main

import (
	"bufio"
	"fmt"
	"os/exec"
	"strings"
)

// Proxy listening information:
//   Address:             127.0.0.1
//   Connection Limit:    0
//   Expiration:          Sat, 02 Mar 2024 21:27:16 CET
//   Port:                63319
//   Protocol:            tcp
//   Session ID:          {id}

//   Credentials:
//     Credential Source ID:   {id}
//     Credential Source Name: {name}
//     Credential Store ID:    {id}
//     Credential Store Type:  vault-generic
//     Secret:
//         {
//               "password": "",
//               "username": ""
//         }

type Credentials struct {
	Target
	username string
	password string
	port     string
}

func Connect_to_target(target Target) (Credentials, error) {
	credentials := Credentials{}
	cmd := exec.Command("boundary", "connect", "-target-id="+target.ID)
	// Get a pipe to read from standard out
	r, _ := cmd.StdoutPipe()

	// Use the same pipe for standard error
	cmd.Stderr = cmd.Stdout

	// Make a new channel which will be used to ensure we get all output
	done := make(chan struct{})

	// Create a scanner which scans r in a line-by-line fashion
	scanner := bufio.NewScanner(r)

	// Use the scanner to scan the output line by line and log it
	// It's running in a goroutine so that it doesn't block
	go func() {

		// Read line by line and process it
		for scanner.Scan() {
			line := scanner.Text()
			if strings.Contains(line, "Port") {
				//parse port
				port := strings.ReplaceAll(strings.ReplaceAll(line, "Port: ", ""), " ", "")
				credentials.port = port
			}
			if strings.Contains(line, "username") {
				//parse username
				copyLine := strings.ReplaceAll(line, "\"", "")
				copyLine = strings.ReplaceAll(copyLine, ":", "")
				copyLine = strings.ReplaceAll(copyLine, ",", "")
				copyLine = strings.ReplaceAll(copyLine, " ", "")
				username := strings.ReplaceAll(copyLine, "username", "")

				credentials.username = username
			}
			if strings.Contains(line, "password") {
				//parse password
				copyLine := strings.ReplaceAll(line, "\"", "")
				copyLine = strings.ReplaceAll(copyLine, ":", "")
				copyLine = strings.ReplaceAll(copyLine, ",", "")
				copyLine = strings.ReplaceAll(copyLine, " ", "")
				password := strings.ReplaceAll(copyLine, "password", "")
				credentials.password = password
			}
			if credentials.port != "" && credentials.username != "" && credentials.password != "" {
				done <- struct{}{}
			}
		}
		done <- struct{}{}
	}()

	// Start the command and check for errors
	err := cmd.Start()
	if err != nil {
		fmt.Println(err)
		return credentials, err
	}

	// Wait for all output to be processed
	<-done

	// Wait for the command to finish
	credentials.Target = target

	return credentials, nil
}
