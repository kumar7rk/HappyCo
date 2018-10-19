package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	intercom "gopkg.in/intercom/intercom-go.v2"
	"io/ioutil"
	"net/http"
	"os"
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
func listSnoozedConversations() /*[]intercom.Conversation*/ {
	p := fmt.Println

	//such a clunky way

	//rohit, richelle, anna, p1, p1a, p2a, p2b, p2c, p3
	var whoseInbox = []string{"rohit", "richelle", "anna", "p1", "p1a", "p2a", "p2b", "p2c", "p3"}
	var allInboxes = []string{"1544605","424979","2687597","1054048","2264830","1615207","931138","2340320","1398520"}
	
	//running all the boxes
	va := 0
	for _, inbox := range allInboxes {
		p(whoseInbox[va])
		va=va+1;
		//running here to get the totalPages of opened conversations for a box
		convoList, err := ic.Conversations.ListByAdmin(&intercom.Admin{ID: json.Number(inbox)}, intercom.SHOW_OPEN, intercom.PageParams{})
		if err != nil {
				fmt.Printf("Error from Intercom listing all opened conversations for: %v\n", err)
		}
		p(convoList.Pages)
	
		//running on all the pages
		for i := 1; i <= int(convoList.Pages.TotalPages); i++ {
	
			convoList, err := ic.Conversations.ListByAdmin(&intercom.Admin{ID: json.Number(inbox)}, intercom.SHOW_OPEN, intercom.PageParams{Page:int64(i)})
			totalConversations := len(convoList.Conversations)
			// fmt.Printf("%s %d %s %d %s","totalConversations in page ",i, ":", totalConversations, "\n")
			if err != nil {
				fmt.Printf("Error from Intercom listing all opened conversations: %v\n", err)
			}
	
			for i := 0; i < totalConversations; i++ {
				convoID := convoList.Conversations[i].ID
				// convo, _ := ic.Conversations.Find(convoID)
				p(convoID)
				// p(convo.ConversationParts)
				// p(convo.Open)
			
				/*l := convo.ConversationParts.Parts
				p("l:")
				p(l[len(l)-2])*/
			}
		}

	}

}

//requirement check if the conversation is currently snoozed.
