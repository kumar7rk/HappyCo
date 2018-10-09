package main

import (
	"strings"
	"time"
)

//********************************************Init********************************************
func init() {
	repCommands["follow"] = RepCommand{Func:followUpConversation , Description: `Follow up message is sent. Conversation is snoozed for 5 days. After 5 days a closing message is sent
	<b>Arguments</b> name
	<b>Want to cancel? You can't. JK. Just enter a note during this time</b>
	It will also cancel if customer messages during the snoozed duration`}
}

//********************************************Sending Buildium reply********************************************
func followUpConversation(user User, author Author, conversationID string, params ...string) {
	userName := "there"

	if user.Email != "" {
		userName = strings.Split(user.Name, " ")[0]
	}
	
	message:=followUpMessage(userName,author.Name)
	addReply(author.ID, conversationID, message)
	addNote(conversationID,"Follow up sent")
	snoozeConversation(conversationID, 3 * 24 * time.Hour)
}
