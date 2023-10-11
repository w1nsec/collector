package config

import (
	"flag"
	"os"
	"strconv"
	"time"
)

const (
	defaultPollInterval   = 2 * time.Second
	defaultReportInterval = 10 * time.Second
)

func AgentSelectArgs(addr *string, pollInterval, reportInterval *int) {
	var (
		flagAddr          string
		flagPoll, flagRep int
	)
	flag.StringVar(&flagAddr, "a", "localhost:8080",
		"address for metric server")
	flag.IntVar(&flagPoll, "p", int(defaultPollInterval.Seconds()),
		"frequency of gathering metrics")
	flag.IntVar(&flagRep, "r", int(defaultReportInterval.Seconds()),
		"frequency of sending metrics")
	flag.Parse()

	if *addr = os.Getenv("ADDRESS"); *addr == "" {
		*addr = flagAddr
	}

	envPoll, err := strconv.Atoi(os.Getenv("POLL_INTERVAL"))
	if err == nil {
		*pollInterval = envPoll
	} else {
		*pollInterval = flagPoll
	}

	envRep, err := strconv.Atoi(os.Getenv("REPORT_INTERVAL"))
	if err == nil {
		*reportInterval = envRep
	} else {
		*reportInterval = flagRep
	}
}
