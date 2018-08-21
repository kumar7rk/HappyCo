package main

import (
	"fmt"
	intercom "gopkg.in/intercom/intercom-go.v2"
)

// structs for reading payload in json received from Intercom
type Author struct {
	Name string `json:"name"`
}

type Part struct {
	Body   string `json:"body"`
	Author Author `json:"author"`
}

type ConversationPart struct {
	Part []Part `json:"conversation_parts"`
}

type ConversationMessage struct {
	Body    string `json:"body"`
	Subject string `json:"subject"`
}

type User struct {
	UserID string `json:"user_id"`
	Type   string `json:"type"`
	Name   string `json:"name"`
	Email  string `json:"email"`
}

type Item struct {
	ConversationID      string              `json:"id"`
	User                User                `json:"user"`
	ConversationMessage ConversationMessage `json:"conversation_message"`
	ConversationPart    ConversationPart    `json:"conversation_parts"`
}

type Data struct {
	Item Item `json:"item"`
}

type Message struct {
	Data Data `json:"data"`
}

//********************************************Add note********************************************
func addNote(conversationID, note string) {
	_, err := ic.Conversations.Reply(conversationID, intercom.Admin{ID: "207278"}, intercom.CONVERSATION_NOTE, note)
	if err != nil {
		fmt.Printf("Error from Intercom while adding note: %v\n", err)
	}
}

//********************************************Reply to user********************************************
func addReply(conversationID, reply string) {
	_, err := ic.Conversations.Reply(conversationID, intercom.Admin{ID: "207278"}, intercom.CONVERSATION_COMMENT, reply)
	if err != nil {
		fmt.Printf("Error from Intercom while adding reply: %v\n", err)
	}
}
