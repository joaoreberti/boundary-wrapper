package main

import (
	"fmt"
	"os"

	"github.com/joho/godotenv"
)

type Target struct {
	ID string
	// ScopeID           string
	// Version           string
	// Type              string
	Name string
	// AuthorizedActions []string
}

func main() {
	err := godotenv.Load()
	if err != nil {
		fmt.Println("Error loading .env file")

	}
	boundaryAddress := os.Getenv("BOUNDARY_ADDRESS")
	if boundaryAddress == "" {
		fmt.Println("Error loading BOUNDARY_ADDRESS")
		return
	}

	dbeaverConfigPath := os.Getenv("DBEAVER_CONFIG_PATH")
	if dbeaverConfigPath == "" {
		fmt.Println("Error loading DBEAVER_CONFIG_PATH")
		return
	}

	// authenticate first
	authenticate_if_needed()

	targets := Get_all_targets()

	if len(targets) == 0 {
		fmt.Println("No targets found")
		return
	}

	// search for targets
	fmt.Print("Search for target: ")

	search, err := Retrieve_user_input()
	if err != nil {
		fmt.Println("Error retrieving input:", err)
		return
	}

	ids, err := Match_targets(targets, search)
	if err != nil {
		fmt.Println("Error matching targets:", err)
		return
	}

	filteredTargets := []Target{}

	for _, target := range targets {
		for _, id := range ids {
			if target.ID == id {
				filteredTargets = append(filteredTargets, target)
			}
		}
	}

	// display the results in the terminal as a list that can be selected
	fmt.Println("Select a target to connect to:")
	for i, filteredTarget := range filteredTargets {
		fmt.Printf("%d) %s\n", i+1, filteredTarget.Name)
	}
	fmt.Println("0) Cancel")
	fmt.Print("Selection: ")
	selection, err := Retrieve_selection_input()
	if err != nil {
		fmt.Println("Error retrieving input:", err)
		return
	}
	if selection == 0 {
		fmt.Println("Cancelled")
		return
	}
	fmt.Println("Connecting to", filteredTargets[selection-1])

	// connect to the selected target
	credentials, err := Connect_to_target(filteredTargets[selection-1])
	if err != nil {
		fmt.Println("Error connecting to target:", err)
	}

	fmt.Println("Enter db name: ")
	dbname, err := Retrieve_user_input()
	if err != nil {
		fmt.Println("Error retrieving input:", err)
		return
	}

	Setup_dbeaver(credentials, dbname)
	fmt.Println("finished script, exiting")
}
