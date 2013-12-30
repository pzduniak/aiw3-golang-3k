package main

import (
	"./config"
	"database/sql"
	"github.com/codegangsta/martini"
	_ "github.com/lib/pq"
	"github.com/robfig/cron"
	"github.com/yvasiyarov/gorelic"
	"log"
	"net/http"
)

var db *sql.DB
var agent *gorelic.Agent
var cfg config.Config

func main() {
	// Load the configuration into a package-scope variable
	cfg = config.Load()

	// Start the NewRelic client if it's enabled in the config file
	if cfg.NewRelic.Enabled {
		agent = gorelic.NewAgent()
		agent.Verbose = cfg.NewRelic.Verbose
		agent.NewrelicName = cfg.NewRelic.Name
		agent.NewrelicLicense = cfg.NewRelic.License
		agent.Run()
	}

	// Connect to the database if logging is enabled in config
	if cfg.Logging.Enabled {
		connection, err := sql.Open("postgres", cfg.Logging.ConnectionString)
		if err != nil {
			log.Printf("Connection to the database failed; %s", err)
		}

		// Export the database connection to package scope
		db = connection
	}

	// Initialize a new martini server
	server := martini.Classic()

	// Allow CORS
	server.Use(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.Header().Set("Access-Control-Allow-Origin", "*")
	})

	// Bind the routes
	server.Get("/total", handleTotal)
	server.Get("/", handleList)

	// Initialize a new scheudler
	scheudle := cron.New()
	scheudle.AddFunc("@every 1m", prepareCache)

	// Call the cache function for the first time
	prepareCache()

	// Run threads
	scheudle.Start()

	// Run SSL if enabled
	if cfg.Binding.HttpsEnabled {
		go func() {
			err := http.ListenAndServeTLS(cfg.Binding.HttpsAddress,
				cfg.Binding.HttpsCertificatePath,
				cfg.Binding.HttpsKeyPath,
				server)

			if err != nil {
				log.Fatalf("Error while serving HTTPS; %s", err)
			}
		}()
	}

	// Run HTTP
	err := http.ListenAndServe(cfg.Binding.HttpAddress, server)
	if err != nil {
		log.Fatalf("Error while serving HTTP; %s", err)
	}
}
