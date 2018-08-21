package main

import (
	"encoding/json"
	"fmt"
	_ "github.com/lib/pq"
	"io/ioutil"
	"net/http"
)

func newAdminNote(w http.ResponseWriter, r *http.Request) {
	fmt.Println("newAdminNote")
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

func processNewAdminNote(user User, conversationID string, note string, author string) {
	if note == "<p>yumi run note</p>" {
		makeAndSendNote(user, conversationID)
	} else if note == "<p>yumi run buildium</p>" {
		sendBuildiumReply(user, conversationID)
	} else if note == "<p>yumi run password</p>" {
		sendPasswordReply(user, conversationID)
	} else if note == "<p>yumi help</p>" || note == "<p>yumi run help</p>" || note == "<p>yumi</p>" {
		listRunCommands(author, conversationID)
	}
}
