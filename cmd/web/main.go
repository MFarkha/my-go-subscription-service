package main

import (
	"database/sql"
	"encoding/gob"
	"fmt"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"

	"log"
	"net/http"

	"github.com/MFarkha/my-go-subscription-service/data"
	"github.com/alexedwards/scs/redisstore"
	"github.com/alexedwards/scs/v2"
	"github.com/gomodule/redigo/redis"
	_ "github.com/jackc/pgconn"
	_ "github.com/jackc/pgx/v5"
	_ "github.com/jackc/pgx/v5/stdlib"
)

const PORT = "3000"

func main() {
	// connect to the database
	db := initDB()

	// create sessions
	session := initSession()

	// create loggers
	infoLog := log.New(os.Stdout, "INFO:\t", log.Ldate|log.Ltime)
	errorLog := log.New(os.Stdout, "ERROR:\t", log.Ldate|log.Ltime|log.Lshortfile)

	// create channels

	// create waitgroup
	wg := sync.WaitGroup{}

	// set up the application config
	app := Config{
		Session:  session,
		DB:       db,
		InfoLog:  infoLog,
		ErrorLog: errorLog,
		Wait:     &wg,
		Models:   data.New(db),
	}
	// set up mail

	// listen for signals (SIGTERM and SIGINT)
	go app.listenForShutdown()

	// listen for web connections
	app.serve()
}

func (app *Config) serve() {
	// start http server
	server := &http.Server{
		Addr:    fmt.Sprintf(":%s", PORT),
		Handler: app.routes(),
	}
	app.InfoLog.Println("Listening on", PORT)

	err := server.ListenAndServe()
	if err != nil {
		app.ErrorLog.Fatalln(err)
	}
}

func initDB() *sql.DB {
	conn := connectToDB()
	if conn == nil {
		log.Panic("cannot connect to database")
	}
	return conn
}

func connectToDB() *sql.DB {
	count := 0
	dsn := os.Getenv("DSN")
	for {
		connection, err := openDB(dsn)
		if err != nil {
			log.Println("database is not ready yet...")
		} else {
			log.Println("connected to database")
			return connection
		}
		count++

		if count > 10 {
			return nil
		}

		log.Println("Backing off for 1 second")
		time.Sleep(1 * time.Second)
	}
}

func openDB(dsn string) (*sql.DB, error) {
	db, err := sql.Open("pgx", dsn)
	if err != nil {
		return nil, err
	}
	err = db.Ping()
	if err != nil {
		return nil, err
	}
	return db, nil
}

func initSession() *scs.SessionManager {
	gob.Register(data.User{})
	// set up sessions
	session := scs.New()
	session.Store = redisstore.New(initRedis())
	session.Lifetime = 24 * time.Hour
	session.Cookie.Persist = true
	session.Cookie.SameSite = http.SameSiteLaxMode
	session.Cookie.Secure = true
	return session
}

func initRedis() *redis.Pool {
	redisPool := &redis.Pool{
		MaxIdle: 10,
		Dial: func() (redis.Conn, error) {
			return redis.Dial("tcp", os.Getenv("REDIS"))
		},
	}
	return redisPool
}

// gracefull shutdown
func (app *Config) listenForShutdown() {
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	app.shutdown()
	os.Exit(0)
}

func (app *Config) shutdown() {
	app.InfoLog.Println("would run clean up tasks...")

	// block until waitgroup is empty
	app.Wait.Wait()

	app.InfoLog.Println("closing channels and shutting down the application")
}
