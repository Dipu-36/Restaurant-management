package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"
)

const version = "1.0.0"

type config struct {
	port int
	env  string
}

type application struct {
	config config
	logger *log.Logger
}

func main() {
	var cfg config

	// flag.InVar registers a command-line flag that binds directly to a variable
	// There are 4 parameters :-
	// 1st param is the pointer to the variable where the parsed value is being stored
	// 2nd param is the flag name used on the CLi like -port=8080
	// 3rd param is thedefault value and the last param is the helper text shown when -help is used
	flag.IntVar(&cfg.port, "port", 4000, "API server port")
	flag.StringVar(&cfg.env, "env", "development", "Envirionment (development | staging | deployment)")
	flag.Parse()

	// the 1st param is the destination where the logs will be written
	// 2nd param is the prefix string will be used further in case of debugging and others
	// 3rd param auto prepends time and date beofre the msg
	logger := log.New(os.Stdout, "", log.Ldate|log.Ltime)

	app := &application{
		config: cfg,
		logger: logger,
	}

	srv := &http.Server{
		Addr:         fmt.Sprintf(":%d", cfg.port),
		Handler:      app.routes(),
		IdleTimeout:  time.Minute,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 10 * time.Second,
	}

	logger.Printf("Starting %s server on %d:", cfg.env, cfg.port)

	err := srv.ListenAndServe()
	logger.Fatal(err)

}
