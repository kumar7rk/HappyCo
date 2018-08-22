/* Adding User Data from DB to Intercom as a note for every new conversation */

package main

import (
	"errors"
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"

	"happyco/apps/tools/libs/buildvars"
)

//********************************************Variable declaration********************************************

var db *sqlx.DB


//********************************************Main function********************************************

func main() {
	var err error

	postgresURI := os.Getenv("POSTGRES_URI")
	if postgresURI == "" {
		fmt.Println("URI error")
	}

	db, err = sqlx.Connect("postgres", postgresURI)

	if err == nil {
		fmt.Println("DB Connected")
	}
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error connecting to DB%v: %v\n", "", err)
		return
	}

	//handling every new convesations in newConversation method in conversation.go
	http.HandleFunc("/conversation", newConversation)
	http.HandleFunc("/note", newAdminNote)
	http.HandleFunc("/healthcheck", healthcheck)
	if err := http.ListenAndServe(":8080", nil); err != nil {
		panic(err)
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
