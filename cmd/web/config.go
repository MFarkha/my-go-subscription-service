package main

import (
	"database/sql"
	"log"
	"sync"

	"github.com/MFarkha/my-go-subscription-service/data"
	"github.com/alexedwards/scs/v2"
)

type Config struct {
	Session  *scs.SessionManager
	DB       *sql.DB
	InfoLog  *log.Logger
	ErrorLog *log.Logger
	Wait     *sync.WaitGroup
	Models   data.Models
}