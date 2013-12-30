package config

type Config struct {
	Master   master_settings
	Binding  binding_settings
	Logging  logging_settings
	NewRelic newrelic_settings `toml:"newrelic"`
}

type master_settings struct {
	Address string
}

type binding_settings struct {
	HttpAddress          string `toml:"http_address"`
	HttpsEnabled         bool   `toml:"https_enabled"`
	HttpsAddress         string `toml:"https_address"`
	HttpsCertificatePath string `toml:"https_certificate_path"`
	HttpsKeyPath         string `toml:"https_key_path"`
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
