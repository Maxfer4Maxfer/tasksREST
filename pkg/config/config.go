package config

import (
	"flag"
	"os"
)

// Config contains all configuration of App
type Config struct {
	DSN      string
	HTTPAddr string
}

// GetConfig returns a fulfilled Config
func GetConfig() Config {
	cfg := Config{}

	fs := flag.NewFlagSet("webApp", flag.ExitOnError)

	cfg.DSN = *fs.String("dsn", "root:root@tcp(mysql:3306)/tasks?charset=utf8&parseTime=True&loc=Local", "Database Source Name")
	cfg.HTTPAddr = *fs.String("http-addr", ":80", "HTTP listen address")

	fs.Parse(os.Args[1:])

	return cfg
}
