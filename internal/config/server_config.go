package config

import (
	"flag"
	"os"
	"strconv"
)

type Args struct {
	Addr     string
	LogLevel string

	// increment 9, FILE_STORAGE
	StoreInterval uint64
	StoragePath   string
	Restore       bool

	// increment 10 DB
	DatabaseURL string
}

func ServerArgsParse(args *Args) {
	args.Addr = os.Getenv("ADDRESS")
	args.LogLevel = os.Getenv("LOGLEVEL")

	// increment 9
	args.StoreInterval, _ = strconv.ParseUint(os.Getenv("STORE_INTERVAL"), 10, 64)
	args.StoragePath = os.Getenv("FILE_STORAGE_PATH")
	args.Restore, _ = strconv.ParseBool(os.Getenv("RESTORE"))

	// increment 10
	args.DatabaseURL = os.Getenv("RESTORE")

	var (
		flagAddr     string
		flagLogLevel string

		// increment 9
		flagStoreInterval uint64
		flagStoragePath   string
		flagRestore       bool

		// increment 10
		flagDatabaseStr string
	)

	flag.StringVar(&flagAddr, "a", "localhost:8080", "address for server")
	flag.StringVar(&flagLogLevel, "l", "info", "log level")

	// increment 9, FILE_STORAGE
	flag.Uint64Var(&flagStoreInterval, "i", 300, "interval in seconds for write store data to file")
	flag.StringVar(&flagStoragePath, "f", "/tmp/metrics-db.json", "file for saving metrics")
	flag.BoolVar(&flagRestore, "r", true, "restore from file-db on startup")

	// increment 10, connect to DB
	flag.StringVar(&flagDatabaseStr, "d", "", "DB connect string")

	flag.Parse()

	if args.Addr == "" {
		args.Addr = flagAddr
	}
	if args.LogLevel == "" {
		args.LogLevel = flagLogLevel
	}

	if args.StoreInterval == 0 {
		args.StoreInterval = flagStoreInterval
	}
	if args.StoragePath == "" {
		args.StoragePath = flagStoragePath
	}
	if !args.Restore {
		args.Restore = flagRestore
	}

	if args.DatabaseURL == "" {
		args.DatabaseURL = flagDatabaseStr
	}

}
