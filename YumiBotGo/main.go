package main

import (
	"errors"
	"net/http"
	"os"
	"time"

	"happyco/apps/tools/libs/buildvars"
	"happyco/libs/log"
	
	"github.com/jmoiron/sqlx"
	
	_ "github.com/lib/pq"
)

var db *sqlx.DB

//********************************************Main function********************************************
func main() {
	var err error

	postgresURI := os.Getenv("POSTGRES_URI")
	if postgresURI == "" {
		log.Error.Println("URI Error")
	}

	db, err = sqlx.Connect("postgres", postgresURI)

	if err == nil {
		log.Info.Println("Connected to database")
	}
	if err != nil {
		log.Error.KV("err", err).Println("Error connecting to database")
		return
	}

	go runFollowupChecker()

	//handling every new convesations in newConversation method in conversation.go
	http.HandleFunc("/conversation", newConversation)
	http.HandleFunc("/note", newAdminNote)
	http.HandleFunc("/healthcheck", healthcheck)
	if err := http.ListenAndServe(":8080", nil); err != nil {
		panic(err)
	}
}

func runFollowupChecker() {
	for {
		// checking if we can follow up || close conversations
		followUpProcess()
		time.Sleep(4 * time.Hour)
	}
}

//********************************************Health check********************************************
func healthcheck(w http.ResponseWriter, r *http.Request) {
	if ping() != nil {
		w.WriteHeader(http.StatusServiceUnavailable)
		w.Write([]byte("Database timeout"))
		return
	}
	w.Write([]byte("OK " + buildvars.BuildScmRevisionShort + " " + buildvars.BuildScmStatus))
}

//********************************************Ping********************************************
func ping() error {
	ping := make(chan error, 0)
	timeout := time.After(10 * time.Second)
	go func() {
		ping <- db.Ping()
	}()
	var err error
	select {
	case <-timeout:
		err = errors.New("Postgres ping timeout error")
	case err = <-ping:
	}
	return err
}
