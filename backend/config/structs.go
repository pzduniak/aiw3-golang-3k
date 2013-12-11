package config

type Config struct {
	Master   master_settings
	Logging  logging_settings
	NewRelic newrelic_settings `toml:"newrelic"`
}

type master_settings struct {
	Address string
}

type logging_settings struct {
	Enabled          bool
	ConnectionString string `toml:"connection_string"`
	TableName        string `toml:"table_name"`
}

type newrelic_settings struct {
	Enabled bool
	Verbose bool
	License string
	Name    string
}
