package main

import (
	"encoding/json"
	"fmt"
	_ "github.com/lib/pq"
	"io/ioutil"
	"net/http"
	"os"
	"strings"
)

//********************************************New Conversation********************************************

//gets intercom token, admin list, reads the payload, and post note as a reply in the conversation
func newConversation(w http.ResponseWriter, r *http.Request) {
	fmt.Println("newConversation")
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
		
		planType := "plan type"
		planTypeRec := getUserPlanType(user.UserID)
		for _, plan := range planTypeRec {
			if plan.Type == "buildium" {
				planType = plan.Type
			}
		}
		// buildium responder
		if planType == "buildium" {
			
			buildiumSupport := strings.Contains(user.Email, "@buildium.com")

			if conversationSubject == "" {
				if !buildiumSupport {
					sendBuildiumReply(user, conversationID)
				}
			} else {
				conversationSubject = strings.ToLower(conversationSubject)

				var autoRepliedMessage bool
				var ignorePhrases = []string{"automatic-reply", "automatic reply", "auto-reply", "auto reply", "out of office", "out-of-office", "automatic"}
				//var ignorePhrases = []string{"auto", "out of office", "out-of-office", "automatic"}

				for _, phrase := range ignorePhrases {
					val := strings.Contains(conversationSubject, phrase)
					if val {
						autoRepliedMessage = true
						break
					}
				}
				if !buildiumSupport && !autoRepliedMessage {
					sendBuildiumReply(user, conversationID)
					return
				}
			}
		}

	}
	// change password autoresponder
	conversationMessage = strings.ToLower(conversationMessage)
	var passwordPhrases = []string{"change password", "change my password", "reset password",
		"reset my password", "pasword is incorrect", "manage password", "manage my password", "forgot password", "forgot my password"}

	var passwordReply bool

	for _, phrase := range passwordPhrases {
		val1 := strings.Contains(conversationMessage, phrase)

		if val1 {
			passwordReply = true
			break
		}
	}
	if passwordReply {
		sendPasswordReply(user, conversationID)
	}
}

//********************************************Getting PlanType for Buidlium auto responder********************************************

func getUserPlanType(ID string) (planTypeRec []Plan) {
	err := db.Select(&planTypeRec, "Select plan_type FROM current_subscriptions WHERE business_id IN (SELECT business_id from business_membership WHERE user_id = $1 AND inactivated_at IS NULL)", ID)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error in plan query %v: %v\n", ID, err)
	}
	return
}
