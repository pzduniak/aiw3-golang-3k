package main

import (
	"database/sql"
	"encoding/json"
	"github.com/codegangsta/martini"
	_ "github.com/lib/pq"
	"github.com/robfig/cron"
	"log"
	"net/http"
	"strconv"
)

var db *sql.DB

var cache = map[string]*ServerQuery{}

func countPlayers() int {
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

func PrepareCache() {
	ips := query_aiw3_master()
	for _, ip := range ips {
		data := server_query(ip)
		if data != nil {
			cache[ip] = data
		} else {
			if cache[ip] != nil {
				delete(cache, ip)
			}
		}
	}

	_, err := db.Query(`INSERT INTO history(time, servers, players) VALUES(NOW(), $1, $2)`, len(cache), countPlayers())
	if err != nil {
		log.Fatal(err)
	}
}

func main() {
	h, err := sql.Open("postgres", "user=aiw3 dbname=aiw3_history")
	if err != nil {
		log.Fatal(err)
	}

	db = h

	server := martini.Classic()

	server.Use(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.Header().Set("Access-Control-Allow-Origin", "*")
	})

	server.Get("/total", func() string {
		json, err := json.MarshalIndent(map[string]int {
			"players": countPlayers(),
			"servers": len(cache),
		}, "", "\t")

		if err != nil {
			return "Error: " + err.Error()
		}

		return string(json)
	})

	server.Get("/", func() string {
		json, err := json.MarshalIndent(cache, "", "\t")

		if err != nil {
			return "Error: " + err.Error()
		}

		return string(json)
	})

	scheudle := cron.New()
	scheudle.AddFunc("@every 1m", PrepareCache)

	PrepareCache()

	scheudle.Start()
	server.Run()
}