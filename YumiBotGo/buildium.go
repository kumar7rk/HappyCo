package main

import (
	"os"
	"fmt"
	"strings"
	intercom "gopkg.in/intercom/intercom-go.v2"
)

func sendBuildiumReply(user User, conversationID string, params... string){
	accessToken := os.Getenv("INTERCOM_ACCESS_TOKEN")
	ic := intercom.NewClient(accessToken, "")
	// if you don't give a name use the name user have in his account
	
	if len(params) == 0 {
		params = strings.Split(user.Name, " ")
	}
	name := "there"
	if len(params) > 0 {
		name = params[0]
	}
	message := "Hi " + name + " 👋 \n \n <b>Our friends at Buildium support your Happy Inspector subscription. \n\n They can be reached at 888-414-1988, or by submitting a ticket through your Buildium account.</b>\n\nBuildium Support team are the best place to help you with this query as they understand your unique workflow and are trained in Happy Inspector 💫 \n \n Please also feel free to take a look through our FAQ on the Buildium integration:  \n https://intercom.help/happyco/frequently-asked-questions/buildium-integration-faq/faq-buildium-integration  \n Thanks!  \n HappyBot ☺"		

	_, err := ic.Conversations.Reply(conversationID, intercom.Admin{ID: "207278"}, intercom.CONVERSATION_COMMENT, message)

		if herr, ok := err.(intercom.IntercomError); ok && herr.GetCode() == "not_found" {
			fmt.Fprintf(os.Stderr, "Error from Intercom when replying to Buildium %v: %v\n", "", err)
			return
		}

}