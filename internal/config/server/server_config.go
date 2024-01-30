package server

import (
	"flag"
	"os"
	"strconv"

	"github.com/rs/zerolog/log"
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

	// increment 14
	Key string

	// increment 21
	CryptoKey string
}

// ServerArgsParse return params in Args struct, that need for server successfully run
func ServerArgsParse() *Args {
	args := new(Args)

	args.Addr = os.Getenv("ADDRESS")
	args.LogLevel = os.Getenv("LOGLEVEL")

	// increment 9
	var err error
	args.StoreInterval, err = strconv.ParseUint(
		os.Getenv("STORE_INTERVAL"), 10, 64)
	if err != nil {
		log.Error().Err(err).Send()
	}
	args.StoragePath = os.Getenv("FILE_STORAGE_PATH")
	args.Restore, err = strconv.ParseBool(os.Getenv("RESTORE"))
	if err != nil {
		log.Error().Err(err).Send()
	}

	// increment 10
	args.DatabaseURL = os.Getenv("DATABASE_DSN")

	// increment 14
	args.Key = os.Getenv("KEY")

	//
	args.CryptoKey = os.Getenv("CryptoKey")

	var (
		flagAddr     string
		flagLogLevel string

		// increment 9
		flagStoreInterval uint64
		flagStoragePath   string
		flagRestore       bool

		// increment 10
		flagDatabaseStr string

		// increment 14
		flagKey string

		// increment 21
		flagCryptoKey string
	)

	flag.StringVar(&flagAddr, "a", "localhost:8080",
		"address for transport")
	flag.StringVar(&flagLogLevel, "l", "info",
		"log level")

	// increment 9, FILE_STORAGE
	flag.Uint64Var(&flagStoreInterval, "i", 300,
		"interval in seconds for write store data to file")
	flag.StringVar(&flagStoragePath, "f", "/tmp/metrics-db.json",
		"file for saving metrics")
	flag.BoolVar(&flagRestore, "r", true,
		"restore from file-db on startup")

	// increment 10, connect to DB
	flag.StringVar(&flagDatabaseStr, "d", "", "DB connect string")

	// increment 14, generate hash for requests body
	flag.StringVar(&flagKey, "k", "", "salt for hmac")

	// increment 21, decrypt requests
	flag.StringVar(&flagCryptoKey, "crypto-key", "",
		"rsa private key path (in pem format), used for encrypt messages")

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

	if args.Key == "" {
		args.Key = flagKey
	}

	if args.CryptoKey == "" {
		args.CryptoKey = flagCryptoKey
	}

	return args
}
