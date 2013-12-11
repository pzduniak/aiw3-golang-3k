package main

import (
	"./aiw3"
	"log"
	"strconv"
)

var cache = map[string]*aiw3.ServerQuery{}

// Populates the cache variable with server query results
func prepareCache() {
	ips := aiw3.QueryMasterServer(cfg.Master.Address)
	for _, ip := range ips {
		data := aiw3.QueryDedicatedServer(ip)
		// Server must return a valid response and have a hostname
		if data != nil && data.Hostname != "" {
			cache[ip] = data
		} else {
			// If the response is invalid, remove the server
			// Invalid response usually means that the server is down
			if cache[ip] != nil {
				delete(cache, ip)
			}
		}
	}

	if cfg.Logging.Enabled {
		_, err := db.Exec(`INSERT INTO `+cfg.Logging.TableName+`(time, servers, players) VALUES(NOW(), $1, $2)`, len(cache), countTotalPlayers())
		if err != nil {
			log.Print(err)
		}
	}
}

// Sums up server.Online of each server in the cache
func countTotalPlayers() int {
	players := 0

	for _, server := range cache {
		num, err := strconv.Atoi(server.Online)
		if err != nil {
			players += 0
		} else {
			players += num
		}
	}

	return players
}