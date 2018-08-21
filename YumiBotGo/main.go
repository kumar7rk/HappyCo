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

	intercom "gopkg.in/intercom/intercom-go.v2"
	"happyco/apps/tools/libs/buildvars"
)

//********************************************Variable declaration********************************************

var db *sqlx.DB

var ic *intercom.Client

func init() {
	accessToken := os.Getenv("INTERCOM_ACCESS_TOKEN")
	ic = intercom.NewClient(accessToken, "")
}

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
	//http.HandleFunc("/", newConversation)
	http.HandleFunc("/", newAdminNote)
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

//**** Intercom ***
func AddNote(conversationID, note string) {
	_, err := ic.Conversations.Reply(conversationID, intercom.Admin{ID: "207278"}, intercom.CONVERSATION_NOTE, note)
	if err != nil {
		fmt.Printf("Error from Intercom while adding note: %v\n", err)
	}
}

func AddReply(conversationID, reply string) {
	_, err := ic.Conversations.Reply(conversationID, intercom.Admin{ID: "207278"}, intercom.CONVERSATION_COMMENT, reply)
	if err != nil {
		fmt.Printf("Error from Intercom while adding reply: %v\n", err)
	}
}
