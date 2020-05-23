package main

import (
	"fmt"
	"net/http"

	"github.com/gorilla/mux"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

// initialize default config values
func restRegistryDefaults() {
	viper.SetDefault("logging.level", log.WarnLevel.String())
	viper.SetDefault("logging.output", "stdout")
	viper.SetDefault("server.host", "localhost")
	viper.SetDefault("server.port", "8000")
	viper.SetDefault("routers.home.mount", "/")
	viper.SetDefault("routers.protected.mount", "protected/")
	viper.SetDefault("routers.unprotected.mount", "unprotected/")
}

func restApp(appName string) {
	configRegistryInit(appName, restRegistryDefaults)
	loggingInit()

	// application mount points
	mountPathHome := viper.GetString("routers.home.mount")
	mountPathProtected := mountPathHome + viper.GetString("routers.protected.mount")
	mountPathUnprotected := mountPathHome + viper.GetString("routers.unprotected.mount")

	router := mux.NewRouter()
	// Order of adding routes is important
	router.PathPrefix(mountPathProtected).HandlerFunc(handlerFuncRESTProtected)
	router.PathPrefix(mountPathUnprotected).HandlerFunc(handlerFuncRESTUnprotected)
	router.PathPrefix(mountPathHome).HandlerFunc(handlerFuncRESTHome)

	log.Println("REST API server started")
	serverURL := viper.GetString("server.host") + ":" + viper.GetString("server.port")
	log.Fatal(http.ListenAndServe(serverURL, router))

}

func handlerFuncRESTHome(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Hello from home page, %s!", r.URL.Path[1:])
}

func handlerFuncRESTProtected(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Hello from protected page, %s!", r.URL.Path[1:])
}

func handlerFuncRESTUnprotected(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Hello from unprotected page, %s!", r.URL.Path[1:])
}
