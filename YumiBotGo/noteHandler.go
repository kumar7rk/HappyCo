package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
	"time"
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
	note- internal message
	author- team member who wrote the note
	*/

	user := msg.Data.Item.User
	conversationId := msg.Data.Item.ConversationID
	note := msg.Data.Item.ConversationPart.Part[0].Body
	author := msg.Data.Item.ConversationPart.Part[0].Author

	go processNewAdminNote(user, author, conversationId, note)
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Received"))
}

//********************************************Processing new note********************************************
func processNewAdminNote(user User, author Author, conversationID string, note string) {
	note = strings.ToLower(note)

	if note == "<p>yumi help</p>" || note == "<p>yumi</p>" {
		listRunCommands(author.Name, conversationID)
		return
	}
	if note == "<p>yumi convo</p>" {
		_ = listOpenedConversations()
		return
	}
	if strings.HasPrefix(note, "<p>yumi rep ") {
		note = strings.TrimSuffix(note[12:], "</p>")
		if strings.Contains(note,"<br>") {
			note = strings.TrimSuffix(note, "<br>")
		}
		params := strings.Split(note, " ")
		if cmd, ok := repCommands[params[0]]; ok {
			cmd.Func(user, author, conversationID, params[1:]...)
		} else {
			fmt.Println("Unable to run", params)
			listRunCommands(author.Name, conversationID)
		}
	}

	if strings.HasPrefix(note, "<p>yumi get ") {
		note = strings.TrimSuffix(note[12:], "</p>")
		if strings.Contains(note,"<br>") {
			note = strings.TrimSuffix(note, "<br>")
		}
		params := strings.Split(note, " ")
		if cmd, ok := getCommands[params[0]]; ok {
			cmd.Func(user, conversationID, params[1:]...)
		} else {
			fmt.Println("Unable to run", params)
			listRunCommands(author.Name, conversationID)
		}
	}

	if note == "<p>###</p>" {
		snoozeConversation(conversationID, 7*24*time.Hour)
	}
}

type RepCommand struct {
	Description string
	Func        func(user User, author Author, conversationID string, params ...string)
}

type Command struct {
	Description string
	Func        func(user User, conversationID string, params ...string)
}

var repCommands map[string]RepCommand = make(map[string]RepCommand)
var getCommands map[string]Command = make(map[string]Command)
