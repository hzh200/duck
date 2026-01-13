package main

import (
	"duck/kernel/log"
	"duck/kernel/manage"
	"duck/kernel/persistence"
	"duck/kernel/server"
	"encoding/json"
	"flag"
	"os"
	"path/filepath"
)

type Config struct {
	devMode bool
	dbPath string
	settingPath string
}

type Setting struct {
	LaunchOnStartup bool
	SlientMode bool
	CloseToTray bool
	KernelPort uint16
	DownloadDirectory string
	Proxy ProxySetting
	TrafficLimit TrafficLimit
}

type ProxySetting struct {
	ProxyMode string
	Host string
	Port int16
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

var config Config
var setting Setting
var manager *manage.Manager

func parseArgs() (Config, error) {
	exe, err := os.Executable()
    if err != nil {
        return Config{}, err
    }

	mode := flag.String("mode", "dev", "kernel running mode: dev/proc")
	dbPath := flag.String("dbPath", filepath.Join(filepath.Dir(exe), "./duck.db"), "database file path")
	settingPath := flag.String("settingPath", filepath.Join(filepath.Dir(exe), "./setting.json"), "setting file path")
	flag.Parse()

	return Config{
		devMode: *mode == "dev", 
		dbPath: *dbPath,
		settingPath: *settingPath,
	}, nil
}

func initialize() error {
	var err error

	// Parse the program config.
	if config, err = parseArgs(); err != nil {
		return err
	}

	// Init the gloal setting.
	var bytes []byte
	if bytes, err = os.ReadFile(config.settingPath); err != nil {
		return err
	}
	json.Unmarshal(bytes, &setting)

	// ticker := time.NewTicker(1 * time.Second)
	// go func() {
	// 	for {
	// 		<-ticker.C
	// 		bytes, err := os.ReadFile(config.settingPath)

	// 		if err != nil {
	// 			log.Error(err)
	// 		}

	// 		json.Unmarshal(bytes, &settings)
	// 	}
	// }()

	// Init the orm.
	var orm *persistence.Persistence
	if orm, err = persistence.InitPersistence(config.dbPath); err != nil {
		return err
	}
	
	// Init the global manager.
	manager = manage.NewManager(orm)
	server.InitRoutes(manager)

	return nil
}

func main() {
	err := initialize()
	if err != nil {
		log.Error(err)
	}

	if config.devMode {
		log.Info("Running under debug mode.")
	}

	// Start the scheduling procedure.
	manager.Schedule()

	// Start the http service.
	err = server.StartServer(setting.KernelPort)
	if err != nil {
		log.Error(err)
	}
}
