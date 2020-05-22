package main

import (
	"fmt"
	"net/http"
	"os"

	"github.com/gorilla/mux"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

// initialize default config values
func configRegistryDefaults() {
	viper.SetDefault("logging.level", log.WarnLevel.String())
	viper.SetDefault("logging.output", "stdout")
	viper.SetDefault("routers.home.mount", "/")
	viper.SetDefault("routers.protected.mount", "protected/")
	viper.SetDefault("routers.unprotected.mount", "unprotected/")
}

func main() {
	configRegistryInit()
	loggingInit()

	mountPathHome := viper.GetString("routers.home.mount")
	mountPathProtected := mountPathHome + viper.GetString("routers.protected.mount")
	mountPathUnprotected := mountPathHome + viper.GetString("routers.unprotected.mount")

	router := mux.NewRouter()
	// Order of adding routes is important
	router.PathPrefix(mountPathProtected).HandlerFunc(handlerFuncProtected)
	router.PathPrefix(mountPathUnprotected).HandlerFunc(handlerFuncUnprotected)
	router.PathPrefix(mountPathHome).HandlerFunc(handlerFuncHome)

	log.Println("http server started")
	log.Fatal(http.ListenAndServe(":8000", router))

}

func handlerFuncHome(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Hello from home page, %s!", r.URL.Path[1:])
}

func handlerFuncProtected(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Hello from protected page, %s!", r.URL.Path[1:])
}

func handlerFuncUnprotected(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Hello from unprotected page, %s!", r.URL.Path[1:])
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
func configRegistryInit() {
	// configure the config
	viper.SetConfigName("config")
	viper.AddConfigPath(".")
	viper.SetConfigType("json")

	configRegistryDefaults()

	// attempt to read the config
	if err := readConfig(); err != nil {
		log.Fatal(fmt.Sprintf("error creating config: %s", err.Error()))
	}
}

// readConfig will attempt to read existing config file, if file doesn't exists it will create one with default settings
func readConfig() error {
	// try to read the config
	err := viper.ReadInConfig()
	if err != nil {
		// write to config if config file does not exists
		err := viper.WriteConfigAs("config.json")
		if err != nil {
			return err
		}
	}
	return nil
}
