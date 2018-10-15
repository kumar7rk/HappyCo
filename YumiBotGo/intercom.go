package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	intercom "gopkg.in/intercom/intercom-go.v2"
	"io/ioutil"
	"net/http"
	"os"
	"reflect"
	"strconv"
	"time"
)

var ic *intercom.Client

func init() {
	accessToken := os.Getenv("INTERCOM_ACCESS_TOKEN")
	ic = intercom.NewClient(accessToken, "")
}

// structs for reading payload in json received from Intercom
type Author struct {
	Name string
	ID   string
}

var YumiBot = Author{Name: "HappyBot", ID: "207278"}

type Part struct {
	Body   string
	Author Author
}

type ConversationPart struct {
	Part []Part `json:"conversation_parts"`
}

type ConversationMessage struct {
	Body    string
	Subject string
}

type User struct {
	UserID string `json:"user_id"`
	Type   string
	Name   string
	Email  string
}

type Item struct {
	ConversationID      string `json:"id"`
	User                User
	ConversationMessage ConversationMessage `json:"conversation_message"`
	ConversationPart    ConversationPart    `json:"conversation_parts"`
}

type Data struct {
	Item Item
}

type Message struct {
	Data Data
}

//********************************************Add note********************************************
func addNote(conversationID, note string) {
	_, err := ic.Conversations.Reply(conversationID, intercom.Admin{ID: "207278"}, intercom.CONVERSATION_NOTE, note)
	if err != nil {
		fmt.Printf("Error from Intercom while adding note: %v\n", err)
	}
}

//********************************************Reply to user********************************************
func addReply(authorID, conversationID string, reply string) {
	_, err := ic.Conversations.Reply(conversationID, intercom.Admin{ID: json.Number(authorID)}, intercom.CONVERSATION_COMMENT, reply)
	if err != nil {
		fmt.Printf("Error from Intercom while adding reply: %v\n", err)
	}
}

//********************************************Assign conversation********************************************
func assignConversation(conversationID string, inboxTo string) {
	_, err := ic.Conversations.Assign(conversationID, &intercom.Admin{ID: "207278"}, &intercom.Admin{ID: json.Number(inboxTo)})
	if err != nil {
		fmt.Printf("Error from Intercom while assigning conversation: %v\n", err)
	}
}

//********************************************Snooze conversation********************************************
func snoozeConversation(conversationID string, duration time.Duration) {
	url := "https://api.intercom.io/conversations/" + conversationID + "/reply"
	payload := []byte(`{ "admin_id":"207278", "message_type":"snoozed", "snoozed_until":` + strconv.FormatInt(time.Now().Add(duration).Unix(), 10) + `}`)

	req, _ := http.NewRequest("POST", url, bytes.NewBuffer(payload))
	req.Header.Set("Authorization", "Bearer "+os.Getenv("INTERCOM_ACCESS_TOKEN"))
	req.Header.Set("Accept", "application/json")
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)

	if err != nil {
		fmt.Printf("Error from Intercom snoozing conversation: %v\n", err)
	}

	defer resp.Body.Close()

	_, _ = ioutil.ReadAll(resp.Body)
}

//********************************************Snooze conversation********************************************
func listAllConversations() {
	// convoList, err := ic.Conversations.ListAll(intercom.PageParams{})

	p := fmt.Println
	convoList, err := ic.Conversations.ListByAdmin(&intercom.Admin{ID: "1544605"}, intercom.SHOW_OPEN, intercom.PageParams{})
	if err != nil {
		fmt.Printf("Error from Intercom listing all opened conversation: %v\n", err)
	}

	p(reflect.TypeOf(convoList))

	/*for _, conv := range convoList {

	}*/

	convoID := convoList.Conversations[0].ID
	convo, _ := ic.Conversations.Find(convoID)
	l := convo.ConversationParts.Parts
	noteContent := l[len(l)-2].Body
	noteAddedTime := l[len(l)-2].CreatedAt

	p("noteAddedTime:")
	p(noteAddedTime) //int64
	fmt.Printf("Current time:")
	p(time.Now().Unix())

	// p("State:"+convoList.Conversations[0].State)

	if l[len(l)-1].PartType == "note" && (noteContent == "yumi rep follow" ||
		noteContent == "<p>yumi convo</p>") {
		p("following up")
		addReply("1544605", "18878341022", followUpMessage("Rohit", "Rohit"))
		addNote("18878341022", "followed up")
	}
	if l[len(l)-1].PartType == "note" && noteContent == "<p>followed up</p>" {
		p("closing conversation")
		addReply("1544605", "18878341022", closingMessage("Rohit", "Rohit"))
	}
}

//requirement check if the conversation is currently snoozed.
