package main

import (
	"fmt"
	"strconv"
	"strings"
	"time"
)

//********************************************Init********************************************
func init() {
	repCommands["follow"] = Command{Func:followUpConversation , Description: `Follow up message is sent. Conversation is snoozed for 3 days. After 3 days a closing message is sent
		Want to cancel? Just enter a note during this time`}
}

//********************************************Sending Buildium reply********************************************
func followUpConversation(user User, conversationID string, params ...string) {
	name := "there"
	snoozeDuration := 7 * 24 * time.Hour

	if user.Email != "" {
		name = strings.Split(user.Name, " ")[0]
	}
	if len(params) == 2 {
		//assuming first param is a name
		name = params[0]
		//checking if the second param is a number or not
		val, err := strconv.Atoi(params[1])
		if err != nil {
			fmt.Printf("Wrong parameters: %v\n", err)
		} else {
			snoozeDuration = time.Duration(val) * 24 * time.Hour
		}
	}
	if len(params) == 1 {
		val, err := strconv.Atoi(params[0])
		if err != nil {
			name = params[0]
		} else {
			snoozeDuration = time.Duration(val) * 24 * time.Hour
		}
	}

	message:=followUpMessage(name,author)
	addReply(conversationID, message)
	assignConversation(conversationID, "1615207")
	snoozeConversation(conversationID, snoozeDuration)
}
