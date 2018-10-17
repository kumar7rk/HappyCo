package main

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"strings"

	_ "github.com/lib/pq"
)

// BuildiumPlanID represents the constant ID for the Buildium plan
const BuildiumPlanID = "8"

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

	var isBuildiumUser bool
	// user.type = lead/user
	if user.Type == "user" {
		makeAndSendNote(user, conversationID)

		planRecs := getUserPlans(user.UserID)
		for _, plan := range planRecs {
			if plan.ID == BuildiumPlanID {
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
					sendBuildiumReply(user, YumiBot, conversationID)
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
					sendBuildiumReply(user, YumiBot, conversationID)
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
		sendPasswordReply(user, YumiBot, conversationID)
	}
	//Sending welcome message for each new conversation including from leads excluding if user is Buildium for is asking to reset password
	userName := "there"
	if user.Email != "" {
		userName = strings.Split(user.Name, " ")[0]
	}

	if !passwordReply && !isBuildiumUser && user.Type == "user" {
		addReply(YumiBot.ID, conversationID, welcomeMessage(userName))
	}
}
