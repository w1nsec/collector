package config

import (
	"flag"
	"os"
)

func ServerArgsParse(addr, logLevel *string) {
	*addr = os.Getenv("ADDRESS")
	*logLevel = os.Getenv("LOGLEVEL")

	var (
		flagAddr     string
		flagLogLevel string
	)
	flag.StringVar(&flagAddr, "a", "localhost:8080", "address for server")
	flag.StringVar(&flagLogLevel, "l", "info", "log level")
	flag.Parse()

	if *addr == "" {
		*addr = flagAddr
	}
	if *logLevel == "" {
		*logLevel = flagLogLevel
	}
}
