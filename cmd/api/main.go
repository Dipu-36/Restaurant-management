package main

import (
	"Dipu-36/restaurant/internals/auth"
	"Dipu-36/restaurant/internals/data"
	"context"
	"database/sql"
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	_ "github.com/lib/pq"
)

const version = "1.0.0"

type config struct {
	port int
	env  string
	db   struct {
		dsn          string
		maxOpenConns int
		maxIdleConns int
		maxIdleTime  string
	}
	jwt struct {
		secret string
	}
}

type application struct {
	config     config
	logger     *log.Logger
	models     data.Models
	db         *sql.DB
	jwtManager *auth.JWTManager
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
	flag.StringVar(&cfg.jwt.secret, "jwt-secret", "change-this-secret-in-production", "JWT signing secret")

	flag.IntVar(&cfg.db.maxOpenConns, "db-max-open-conns", 25, "PostreSQL max open connections")
	flag.IntVar(&cfg.db.maxIdleConns, "db-max-idle-conns", 25, "PostgreSQL max idle connections")
	flag.StringVar(&cfg.db.maxIdleTime, "db-maxidle-time", "15m", "PostgreSQL max connection idle time")
	flag.StringVar(&cfg.db.dsn, "db-dsn", os.Getenv("RESTAURANT_DB_DSN"), "PostgreSql DSN")

	flag.Parse()

	// the 1st param is the destination where the logs will be written
	// 2nd param is the prefix string will be used further in case of debugging and others
	// 3rd param auto prepends time and date beofre the msg
	logger := log.New(os.Stdout, "", log.Ldate|log.Ltime)

	db, err := openDB(cfg)
	if err != nil {
		logger.Fatal(err)
	}

	defer db.Close()

	logger.Printf("database connection pool established")

	app := &application{
		config: cfg,
		logger: logger,
		models: data.NewModels(db),
		db:     db,
		jwtManager: &auth.JWTManager{
			SecretKey: []byte(cfg.jwt.secret),
			Issuer:    "restaurant-api",
			TTL:       24 * time.Hour,
		},
	}

	srv := &http.Server{
		Addr:         fmt.Sprintf(":%d", cfg.port),
		Handler:      app.routes(),
		IdleTimeout:  time.Minute,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 10 * time.Second,
	}

	logger.Printf("Starting %s server on %d:", cfg.env, cfg.port)

	go func() {
		err := srv.ListenAndServe()
		if err != nil && err != http.ErrServerClosed {
			logger.Fatal(err)
		}
	}()

	quit := make(chan os.Signal, 1)

	signal.Notify(
		quit,
		os.Interrupt,
		syscall.SIGTERM,
	)

	<-quit

	logger.Println("Shutting down server...")

	ctx, cancel := context.WithTimeout(
		context.Background(),
		30*time.Second,
	)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		logger.Fatal(err)
	}

	logger.Println("Closing database connection...")

	if err := db.Close(); err != nil {
		logger.Fatal(err)
	}

	logger.Println("Server stopped gracefully.")
}

// openDB() helper returns a sql.DB connection pool
func openDB(cfg config) (*sql.DB, error) {
	// Use sq.Open() to create an empty conection pool using the DSN from the config struct
	db, err := sql.Open("postgres", cfg.db.dsn)
	if err != nil {
		return nil, err
	}

	// Set the maximum number of open connections in pool passing a value
	// less than or equal than to 0 will mean there is no limit
	db.SetMaxOpenConns(cfg.db.maxOpenConns)

	// set themaximum number of idle connections in the pool Again passing a value
	// less than or equal to 0 will mean there is no limit
	db.SetMaxIdleConns(cfg.db.maxIdleConns)

	// Use the time.ParseDuration() function to convert the idle timeout duration
	// string to a time.Duration type
	duration, err := time.ParseDuration(cfg.db.maxIdleTime)
	if err != nil {
		return nil, err
	}

	// Set the maximum idle timeout
	db.SetConnMaxIdleTime(duration)

	// Using PingContext() to establish a new connection to the Database passing in the conetext we
	// created above as a parameter, if the connection couldn't be established succesfully within
	// the 5 second timeout deadline then this will return an error
	err = db.PingContext(ctx)
	if err != nil {
		return nil, err
	}

	return db, nil
}
