package main

import (
	"strings"
	"time"
)

//********************************************Init********************************************
func init() {
	repCommands["follow"] = RepCommand{Func:followUpConversation , Description: `A conversation is snoozed for 3 days. 
	After 2 days a follow up message is sent from you. (a note "follow up sent" is added)
	Again, the conversation is snoozed for a week.
	After 4 days a closing message is sent.
	<b>Want to cancel? You can't. JK. Just enter a note during this time</b>
	It will also cancel if a customer messages or we reply.`}
}

//********************************************Follow up********************************************
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
