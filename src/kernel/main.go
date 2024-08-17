package main

import (
	"duck/kernel/log"
	"duck/kernel/persistence"
	"duck/kernel/server"
	"encoding/json"
	"flag"
	"os"
	"path/filepath"
	"time"
)

type Config struct {
	devMode bool
	port uint16
	dbPath string
	settingPath string
}

type Setting struct {
	LaunchOnStartup bool
	CloseToTray bool
	DownloadDirectory string
	Proxy ProxySetting
	TrafficLimit TrafficLimit
}

type ProxySetting struct {
	ProxyMode string
	Host string
	Port string
}

type ProxyMode string
const (
	Off ProxyMode = "off"
	System ProxyMode = "system"
	Manually ProxyMode = "manually"
)

type TrafficLimit struct {
	Enabled string
	Limit string
}

func main() {
	config, err := parseArgs()

	if err != nil {
		log.Error(err)
	}

	if config.devMode {
		log.Info("Running under debug mode.")
	}

	ticker := time.NewTicker(1 * time.Second)
	go func() {
		for {
			<-ticker.C
			bytes, err := os.ReadFile(config.settingPath)

			if err != nil {
				log.Error(err)
			}
			settings := Setting{}
			json.Unmarshal(bytes, &settings)
		}
	}()

	// Init the orm.
	persistence, err := persistence.InitPersistence(config.dbPath)

	if err != nil {
		log.Error(err)
	}

	// Init the http service.
	err = server.StartServer(config.port, persistence)

	if err != nil {
		log.Error(err)
	}
}

func parseArgs() (Config, error) {
	exe, err := os.Executable()

    if err != nil {
        return Config{}, err
    }

    defaultDBPath := filepath.Join(filepath.Dir(exe), "./duck.db")
	defaultSettingPath := filepath.Join(filepath.Dir(exe), "./setting.json")

	mode := flag.String("mode", "proc", "kernel running mode: dev/proc")
	port := flag.Int("port", 9000, "backend http server listening port")
	dbPath := flag.String("dbPath", defaultDBPath, "database file path")
	settingPath := flag.String("settingPath", defaultSettingPath, "setting file path")
	flag.Parse()

	return Config{
		devMode: *mode == "dev", 
		port: uint16(*port), 
		dbPath: *dbPath,
		settingPath: *settingPath,
	}, nil
}
