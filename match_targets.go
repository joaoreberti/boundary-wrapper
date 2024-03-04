package main

import (
	"strings"
)

func Match_targets(targets []Target, search string) ([]string, error) {
	var matches []string
	for _, target := range targets {
		if strings.Contains(target.Name, search) {
			matches = append(matches, target.ID)
		}
	}
	return matches, nil
}
