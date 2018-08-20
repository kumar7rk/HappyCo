package main

import (
	"encoding/json"
	"fmt"
	_ "github.com/lib/pq"
	"io/ioutil"
	"net/http"
	"os"
	"strings"
)

func newAdminNote(w http.ResponseWriter, r *http.Request) {
	fmt.Println("newConversation")
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
	conversationmessage- user's message
	*/

	user := msg.Data.Item.User
	conversationId := msg.Data.Item.ConversationID
	conversationMessage := msg.Data.Item.ConversationMessage.Body
	conversationSubject := msg.Data.Item.ConversationMessage.Subject

	go processNewConversation(user, conversationId, conversationMessage, conversationSubject)
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Received"))
}
