package main

import (
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/guths/zpe/config"
	"github.com/guths/zpe/handlers"
	"github.com/guths/zpe/router"
)

func init() {
	config.InitializeAppConfig()
}

func main() {
	//initialize all models and handlers
	//the db are instantiated here and passed to all orm models
	if err := handlers.InitializeHandler(); err != nil {
		log.Fatalln(err)
	}

	s := &http.Server{
		Addr: fmt.Sprintf(":%d", config.AppConfig.Port),
		//the router return is a gin engine with all routes and middlewares
		Handler:        router.InitializeRouter(),
		ReadTimeout:    10 * time.Second,
		WriteTimeout:   10 * time.Second,
		MaxHeaderBytes: 1 << 20,
	}

	if err := s.ListenAndServe(); err != nil {
		log.Fatalln(err)
	}
}
