package main

import (
	"fmt"
	"os"

	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

var (
	configPath = "."
	configName = "Config"
	configType = "json"
)

func main() {
	appName := "web"
	switch appName {
	case "webapp":
		webApp(appName)
	case "rest":
		grpcApp(appName)
	}

}

// init should be moved to the application init
func loggingInit() {
	logOutput := viper.GetString("logging.output")

	// Log as JSON instead of the default ASCII formatter.
	log.SetFormatter(&log.JSONFormatter{})

	// Output to stdout instead of the default stderr
	// Can be any io.Writer, see below for File example
	switch logOutput {
	case "stdout":
		log.SetOutput(os.Stdout)
	default:
		fmt.Printf("Invalid config value for logging.output: %s", logOutput)
		os.Exit(1)
	}

	logLevel, _ := log.ParseLevel(viper.GetString("logging.level"))
	log.SetLevel(logLevel)
}

// configInit initilises the Viper config registry
func configRegistryInit(appName string, registryDefaults func()) {
	// configure the config
	viper.SetConfigName(appName + configName)
	viper.AddConfigPath(configPath)
	viper.SetConfigType(configType)

	registryDefaults()

	// attempt to read the config
	if err := readConfig(appName); err != nil {
		log.Fatal(fmt.Sprintf("error creating config: %s", err.Error()))
	}
}

// readConfig will attempt to read existing config file, if file doesn't exists it will create one with default settings
func readConfig(appType string) error {
	// try to read the config
	err := viper.ReadInConfig()
	if err != nil {
		configFile := appType + configName + "." + configType
		// write to config if config file does not exists
		err := viper.WriteConfigAs(configFile)
		if err != nil {
			return err
		}
	}
	return nil
}
