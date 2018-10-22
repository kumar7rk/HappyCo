package main

import (
	"time"
)

//********************************************Init********************************************
func init() {
	repCommands["follow"] = RepCommand{Func: followUpConversation, Description: `A conversation is snoozed for 3 days. 
	After 2 days a follow up message is sent from you. (a note "followed up" is added from HappyBot)
	The conversation is snoozed for 5 days.
	After 4 days a closing message is sent.
	<b>Want to cancel? You can't. JK. Just enter a note during this time</b>
	It will also cancel if a customer messages or we reply.
	Currently, this job runs every 12 hours.`}
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

		canFollowUp := false
		if followUpTime > noteAddedTime {
			canFollowUp = true
		}

		canClose := false
		if closeTime > noteAddedTime {
			canClose = true
		}
		if !canFollowUp && !canClose {
			continue
		}
		authorID := lastNote.Author.ID
		authorName := m[authorID]
		userName := getUserName(convo.User.ID)

		//checking if a conversation meets the parameters for following up
		if lastNote.PartType == "note" && lastNote.Body == "<p>yumi rep follow</p>" && canFollowUp {
			addReply(authorID, convo.ID, followUpMessage(userName, authorName))
			snoozeConversation(convo.ID, 7*24*time.Hour)
			addNote(convo.ID, "followed up")
		}
		//check if conversation can be closed
		if lastNote.PartType == "note" && lastNote.Body == "<p>followed up</p>" && canClose {
			addReply(authorID, convo.ID, closingMessage(userName, authorName))
			addNote(convo.ID, "Conversation closed")
		}
	}
}
