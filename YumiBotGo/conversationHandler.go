package main

import (
	"encoding/json"
	_ "github.com/lib/pq"
	"io/ioutil"
	"net/http"
	"strings"
)

//********************************************New Conversation********************************************
func newConversation(w http.ResponseWriter, r *http.Request) {
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
	conversationmessage- user's message
	*/

	user := msg.Data.Item.User
	conversationId := msg.Data.Item.ConversationID
	conversationMessage := msg.Data.Item.ConversationMessage.Body
	conversationSubject := msg.Data.Item.ConversationMessage.Subject

	go processNewConversation(user, conversationId, conversationMessage, conversationSubject)
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Received"))
}

//********************************************Checking for different attributes********************************************
func processNewConversation(user User, conversationID string, conversationMessage string, conversationSubject string) {
	// user.type = lead/user
	if user.Type == "user" {
		makeAndSendNote(user, conversationID)

		var isBuildiumUser bool
		planTypeRec := getUserPlanType(user.UserID)
		for _, plan := range planTypeRec {
			if plan.Type == "buildium" {
				isBuildiumUser = true
				break
			}
		}
		// buildium autoresponder
		if isBuildiumUser {
			buildiumSupport := strings.Contains(user.Email, "@buildium.com")
			happyCoTeam := strings.Contains(user.Email, "@happy.co")
			if conversationSubject == "" {
				if !buildiumSupport && !happyCoTeam {
					sendBuildiumReply(user, conversationID)
				}
			} else {
				conversationSubject = strings.ToLower(conversationSubject)

				var autoRepliedMessage bool
				var ignorePhrases = []string{"auto", "out of office", "out-of-office"}

				for _, phrase := range ignorePhrases {
					if strings.Contains(conversationSubject, phrase) {
						autoRepliedMessage = true
						break
					}
				}
				if !buildiumSupport && !happyCoTeam && !autoRepliedMessage {
					sendBuildiumReply(user, conversationID)
					return
				}
			}
		}
	}
	// change password autoresponder
	conversationMessage = strings.ToLower(conversationMessage)

	var passwordPhrases = []string{"change password", "change my password", "reset password",
		"reset my password", "resetting my password", "password is incorrect", "password was incorrect", "manage password", "manage my password", "forgot password", "forgot my password"}

	var passwordReply bool

	for _, phrase := range passwordPhrases {
		if strings.Contains(conversationMessage, phrase) {
			passwordReply = true
			break
		}
	}
	//Users replying to password reset email
	text := "Someone has requested a link to change your password"
	if strings.Contains(conversationMessage, text) {
		passwordReply = false
	}

	if passwordReply {
		sendPasswordReply(user, conversationID)
	}
}
