package main

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"os"
	"strconv"
	"time"

	"happyco/libs/log"

	intercom "gopkg.in/intercom/intercom-go.v2"
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

var yumiBot = Author{Name: "HappyBot", ID: "2708121"}

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
	_, err := ic.Conversations.Reply(conversationID, intercom.Admin{ID: json.Number(yumiBot.ID)}, intercom.CONVERSATION_NOTE, note)
	if err != nil {
		log.Error.KV("err", err).KV("conversationID", conversationID).Println("Error from Intercom while adding note")
	}
}

//********************************************Reply to user********************************************
func addReply(authorID, conversationID string, reply string) {
	_, err := ic.Conversations.Reply(conversationID, intercom.Admin{ID: json.Number(authorID)}, intercom.CONVERSATION_COMMENT, reply)
	if err != nil {
		log.Error.KV("err", err).KV("conversationID", conversationID).Println("Error from Intercom while adding reply")
	}
}

//********************************************Assign conversation********************************************
func assignConversation(conversationID string, inboxTo string) {
	_, err := ic.Conversations.Assign(conversationID, &intercom.Admin{ID: json.Number(yumiBot.ID)}, &intercom.Admin{ID: json.Number(inboxTo)})
	if err != nil {
		log.Error.KV("err", err).KV("conversationID", conversationID).KV("inboxTo", inboxTo).Println("Error from Intercom while assigning conversation")
	}
}

//********************************************Snooze conversation********************************************
func snoozeConversation(conversationID string, duration time.Duration) {
	url := "https://api.intercom.io/conversations/" + conversationID + "/reply"
	payload := []byte(`{ "admin_id":`+yumiBot.ID+`, "message_type":"snoozed", "snoozed_until":` + strconv.FormatInt(time.Now().Add(duration).Unix(), 10) + `}`)

	req, _ := http.NewRequest("POST", url, bytes.NewBuffer(payload))
	req.Header.Set("Authorization", "Bearer "+os.Getenv("INTERCOM_ACCESS_TOKEN"))
	req.Header.Set("Accept", "application/json")
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)

	if err != nil {
		log.Error.KV("err", err).KV("conversationID", conversationID).KV("SnoozeDuration", duration).Println("Error from Intercom snoozing conversation")
	}

	defer resp.Body.Close()

	_, _ = ioutil.ReadAll(resp.Body)
}

//********************************************List Opened conversation********************************************
func listOpenedConversations() []intercom.Conversation {

	var allOpenedConversations = []intercom.Conversation{}
	//rohit, richelle, anna, p1, p1a, p2a, p2b, p2c, p3
	var allInboxes = []string{"1544605", "424979", "2687597", "1054048", "2264830", "1615207", "931138", "2340320", "1398520"}

	//running for all the boxes
	for _, inbox := range allInboxes {
		//running here to get the totalPages of opened conversations for a box
		convoList, err := ic.Conversations.ListByAdmin(&intercom.Admin{ID: json.Number(inbox)}, intercom.SHOW_OPEN, intercom.PageParams{})
		if err != nil {
			continue
			log.Error.KV("err", err).KV("inbox", inbox).Println("Error from Intercom listing all opened conversations for an inbox")
		}

		//running on all the pages
		for i := 1; i <= int(convoList.Pages.TotalPages); i++ {
			convoList, err := ic.Conversations.ListByAdmin(&intercom.Admin{ID: json.Number(inbox)}, intercom.SHOW_OPEN, intercom.PageParams{Page: int64(i)})
			if err != nil {
				continue
				log.Error.KV("err", err).KV("inbox", inbox).KV("Page", i).Println("Error from Intercom listing all opened conversations for a page in an inbox")
			}
			totalConversations := len(convoList.Conversations)

			// for all conversations
			for j := 0; j < totalConversations; j++ {
				convoID := convoList.Conversations[j].ID
				convo, err := ic.Conversations.Find(convoID)
				if err != nil {
					log.Error.KV("err", err).KV("conversationID", convoID).Println("Error from Intercom finding a conversation")
					continue
				}
				allOpenedConversations = append(allOpenedConversations, convo)
			}
		}
	}
	return allOpenedConversations
}

//********************************************Getting user name********************************************

func getUserName(ID string) string {
	user, err := ic.Users.FindByID(ID)
	if err != nil {
		log.Error.KV("err", err).KV("intercomUserID", ID).Println("Error from Intercom finding a user")
	}
	return user.Name
}
