package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
)

//********************************************New adming note********************************************
func newAdminNote(w http.ResponseWriter, r *http.Request) {
	// Read body/payload
	b, err := ioutil.ReadAll(r.Body)
	defer r.Body.Close()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Unmarshal the json
	var msg Message
	err = json.Unmarshal(b, &msg)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	/* getting attributes from the received json
	user- user's attributes - name id email type
	conversationId - Intercom conversation ID
	note- team's message
	author- team member who wrote the note
	*/

	user := msg.Data.Item.User
	conversationId := msg.Data.Item.ConversationID
	note := msg.Data.Item.ConversationPart.Part[0].Body
	author := msg.Data.Item.ConversationPart.Part[0].Author.Name

	go processNewAdminNote(user, conversationId, note, author)
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Received"))
}

//********************************************Processing new note********************************************
func processNewAdminNote(user User, conversationID string, note string, author string) {
	if note == "<p>yumi help</p>" || note == "<p>yumi run help</p>" || note == "<p>yumi</p>" {
		listRunCommands(author, conversationID)
		return
	}

	if strings.HasPrefix(note, "<p>yumi run ") {
		note = strings.TrimSuffix(note[12:], "</p>")
		params := strings.Split(note, " ")
		if cmd, ok := commands[params[0]]; ok {
			cmd(user, conversationID, params[1:]...)
		} else {
			fmt.Println("Unable to run", params)
			listRunCommands(author, conversationID)
		}
	}
}

type CommandFunc func(user User, conversationID string, params ...string)

var commands map[string]CommandFunc = make(map[string]CommandFunc)
