package main

import (
	"strings"
	// intercom "gopkg.in/intercom/intercom-go.v2"
)

//********************************************Sending Buildium reply********************************************
func sendBuildiumReply(user User, conversationID string, params ...string) {
	// if you don't give a name use the name user have in his account

	if len(params) == 0 {
		params = strings.Split(user.Name, " ")
	}
	name := "there"
	if len(params) > 0 {
		name = params[0]
	}
	message := "Hi " + name + " 👋 \n \n <b>Our friends at Buildium support your Happy Inspector subscription. \n\n They can be reached at 888-414-1988, or by submitting a ticket through your Buildium account.</b>\n\nBuildium Support team are the best place to help you with this query as they understand your unique workflow and are trained in Happy Inspector 💫 \n \n Please also feel free to take a look through our FAQ on the Buildium integration:  \n https://intercom.help/happyco/frequently-asked-questions/buildium-integration-faq/faq-buildium-integration  \n Thanks!  \n HappyBot ☺"

	addReply(conversationID, message)
	// _, err := ic.Conversations.Assign(conversationID, intercom.Admin{ID: "207278"}, intercom.Admin{ID: "1615207"}) //giving an error *Intercom.admin type...

}
