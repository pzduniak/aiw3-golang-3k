package main

import (
	"encoding/json"
	"github.com/codegangsta/martini"
	"github.com/robfig/cron"
	"net/http"
	"strconv"
)

var cache = map[string]*ServerQuery{}

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
}

func main() {
	server := martini.Classic()

	server.Use(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.Header().Set("Access-Control-Allow-Origin", "*")
	})

	server.Get("/total", func() string {
		players := 0

		for _, server := range cache {
			num, err := strconv.Atoi(server.Online)
			if err != nil {
				players += 0
			} else {
				players += num
			}
		}

		json, err := json.MarshalIndent(map[string]int {
			"players": players,
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