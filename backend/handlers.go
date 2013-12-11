package main

import (
	"encoding/json"
)

// GET / endpoint
func handleList() string {
	json, err := json.Marshal(cache)

	if err != nil {
		return "Error: " + err.Error()
	}

	return string(json)
}

// GET /total endpoint
func handleTotal() string {
	json, err := json.Marshal(map[string]int{
		"players": countTotalPlayers(),
		"servers": len(cache),
	})

	if err != nil {
		return "Error: " + err.Error()
	}

	return string(json)
}
