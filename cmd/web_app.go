package main

import (
	"fmt"
	"net/http"

	"github.com/gorilla/mux"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

// initialize default config values
func webRegistryDefaults() {
	viper.SetDefault("logging.level", log.WarnLevel.String())
	viper.SetDefault("logging.output", "stdout")
	viper.SetDefault("server.host", "localhost")
	viper.SetDefault("server.port", "8000")
	viper.SetDefault("routers.home.mount", "/")
	viper.SetDefault("routers.protected.mount", "protected/")
	viper.SetDefault("routers.unprotected.mount", "unprotected/")
}

func webApp(appName string) {
	configRegistryInit(appName, webRegistryDefaults)
	loggingInit()

	// application mount points
	mountPathHome := viper.GetString("routers.home.mount")
	mountPathProtected := mountPathHome + viper.GetString("routers.protected.mount")
	mountPathUnprotected := mountPathHome + viper.GetString("routers.unprotected.mount")

	router := mux.NewRouter()
	// Order of adding routes is important
	router.PathPrefix(mountPathProtected).HandlerFunc(handlerFuncProtected)
	router.PathPrefix(mountPathUnprotected).HandlerFunc(handlerFuncUnprotected)
	router.PathPrefix(mountPathHome).HandlerFunc(handlerFuncHome)

	log.Println("http server started")
	serverURL := viper.GetString("server.host") + ":" + viper.GetString("server.port")
	log.Fatal(http.ListenAndServe(serverURL, router))

}

func handlerFuncHome(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Hello from home page, %s!", r.URL.Path[1:])
	// Call gRPC uService for data model
}

func handlerFuncProtected(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Hello from protected page, %s!", r.URL.Path[1:])
	// Call gRPC uService for data model
}

func handlerFuncUnprotected(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Hello from unprotected page, %s!", r.URL.Path[1:])
	// Call gRPC uService for data model
}
