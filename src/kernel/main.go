package main

import (
	"duck/kernel/server"
	"duck/kernel/persistence"
	"flag"
	"fmt"
	"os"
    "path/filepath"
)

type Config struct {
	devMode bool
	port uint16
	dbPath string
}

func main() {
	config, err := parseArgs()

	if err != nil {
		
	}

	if config.devMode {
		fmt.Println("Running under development mode.")
	}

	conn, err := persistence.InitConnection(config.dbPath)

	if err != nil {
		
	}

	conn.Close()

	err = server.Run(config.port)

	if err != nil {
		
	}
}

func parseArgs() (Config, error) {
	exe, err := os.Executable()

    if err != nil {
        return Config{}, err
    }

    defaultDBPath := filepath.Join(filepath.Dir(exe), "./duck.db")

	mode := flag.String("mode", "proc", "kernel running mode: dev/proc")
	port := flag.Int("port", 9000, "backend http server listening port")
	dbPath := flag.String("dbPath", defaultDBPath, "database file path")
	flag.Parse()

	return Config{
		devMode: *mode == "dev", 
		port: uint16(*port), 
		dbPath: *dbPath,
	}, nil
}
