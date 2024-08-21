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

func main() {
	config, err := parseArgs()

	if err != nil {
		log.Error(err)
	}

	if config.devMode {
		log.Info("Running under debug mode.")
	}

	settings := Setting{}
	bytes, err := os.ReadFile(config.settingPath)
	if err != nil {
		log.Error(err)
	}
	json.Unmarshal(bytes, &settings)

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
	persistence, err := persistence.InitPersistence(config.dbPath)
	if err != nil {
		log.Error(err)
	}
	
	// Init the main.
	manager := manage.NewManager(persistence)

	// Start the scheduling procedure.
	manager.Schedule()

	// Init the http service.
	server.InitRoutes(manager)
	err = server.StartServer(settings.KernelPort)

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
	dbPath := flag.String("dbPath", defaultDBPath, "database file path")
	settingPath := flag.String("settingPath", defaultSettingPath, "setting file path")
	flag.Parse()

	return Config{
		devMode: *mode == "dev", 
		dbPath: *dbPath,
		settingPath: *settingPath,
	}, nil
}
