package main

import (
	"encoding/json"
	"github.com/codegangsta/martini"
	"github.com/robfig/cron"
)

var cache = map[string]*ServerQuery{}

func PrepareCache() {
	ips := query_aiw3_master()
	for _, ip := range ips {
		data := server_query(ip)
		if data != nil {
			cache[ip] = data
		}
	}
}

func main() {
	server := martini.Classic()
	server.Get("/", func() string {
		json, err := json.Marshal(cache)

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