package main

import (
	"os"
	"fmt"
	intercom "gopkg.in/intercom/intercom-go.v2"
)

func sendPasswordReply(name string, conversationID string){
	accessToken := os.Getenv("INTERCOM_ACCESS_TOKEN")
	ic := intercom.NewClient(accessToken, "")
	
	message := "Hi " + name + " 👋 \n \n You can reset your password either through the web application or iOS application by clicking Forgot Your Password. \n\nAlso be sure to check out our FAQ on managing your profile: \n\n <a href='http://intercom.help/happyco/frequently-asked-questions/happy-inspector-faq-setup-and-user-management/faq-managing-passwords-user-details-and-access-to-your-properties'> FAQ: Managing Passwords, User Details and Access to Your Properties</a> \n\n Let me know if you\n\n still have issues logging in, \n Thanks!  \n HappyCo team ☺"		
	_, err := ic.Conversations.Reply(conversationID, intercom.Admin{ID: "207278"}, intercom.CONVERSATION_COMMENT, message)

		if herr, ok := err.(intercom.IntercomError); ok && herr.GetCode() == "not_found" {
			fmt.Fprintf(os.Stderr, "Error from Intercom when replying to Buildium %v: %v\n", "", err)
			return
		}

}