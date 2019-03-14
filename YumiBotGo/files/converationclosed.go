package main

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
)

func conversationClosed(w http.ResponseWriter, r *http.Request) {
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

	user := msg.Data.Item.User
	conversationID := msg.Data.Item.ConversationID
	author := msg.Data.Item.ConversationPart.Part[0].Author

	go processClosedConversation(user, conversationID, author)
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Received"))

}

func processClosedConversation(user User, conversationID string, author Author) {
	userEmail := user.Email
	businessRec := getBusiness(userEmail)

	if len(businessRec) != 1 {
		return
	}
	for _, business := range businessRec {
		userRoleID := business.Role.String
		var businessPermission = business.PermissionsModel

		if (userRoleID == "1" || userRoleID == "8") && businessPermission == "basic-roles" {
			// addReply(author.ID, conversationID, "Hey you're on a basic plan. Want more control of your account. Hit us up.")
		}
	}
}
