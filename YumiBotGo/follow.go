package main

import (
	"strings"
	"time"
)

//********************************************Init********************************************
func init() {
	replyCommands["follow"] = ReplyCommand{Func: followUpConversation, Description: `A conversation is snoozed for 3 days. 
	After 2 days a follow up message is sent from you. (a note "followed up" is added from HappyBot)
	The conversation is snoozed for a further 5 days.
	After 4 days a closing message is sent.
	<b>Want to cancel? You can't. JK. Just enter a note during this time</b>
	It will also cancel if a customer messages or we reply.`}
}

//********************************************Follow up********************************************
func followUpConversation(user User, author Author, conversationID string, params ...string) {
	snoozeConversation(conversationID, 3*24*time.Hour)
}

func followUpProcess() {
	allOpenedConversations := listOpenedConversations()

	for _, convo := range allOpenedConversations {
		notes := convo.ConversationParts.Parts
		lastNote := notes[len(notes)-1]

		m := make(map[string]string)
		m["1544605"] = "Rohit"
		m["424979"] = "Richelle"
		m["2687597"] = "Anna"

		followUpDuration := 2 * 24 * time.Hour
		closeDuration := 4 * 24 * time.Hour

		noteAddedTime := lastNote.CreatedAt
		followUpTime := time.Now().Add(-followUpDuration).Unix()
		closeTime := time.Now().Add(-closeDuration).Unix()

		canFollowUp := followUpTime > noteAddedTime
		canClose := closeTime > noteAddedTime

		if !canFollowUp && !canClose {
			continue
		}
		authorID := lastNote.Author.ID
		authorName := m[authorID]
		userName := getUserName(convo.User.ID)

		//checking if a conversation meets the parameters for following up
		if lastNote.PartType == "note" && lastNote.Body == "<p>yumi reply follow</p>" && canFollowUp {
			addReply(authorID, convo.ID, followUpMessage(userName, authorName))
			snoozeConversation(convo.ID, 5*24*time.Hour)
			addNote(convo.ID, "Followed up")
		}
		//check if conversation can be closed
		if lastNote.PartType == "note" && strings.EqualFold(lastNote.Body, "<p>Followed up</p>") && canClose {
			addReply(authorID, convo.ID, closingMessage(userName, authorName))
			addNote(convo.ID, "Conversation closed")
		}
	}
}
