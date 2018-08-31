package main

import (
	"strings"
)

//********************************************Init********************************************
func init() {
	commands["buildium"] = sendBuildiumReply
}

//********************************************Sending Buildium reply********************************************
func sendBuildiumReply(user User, conversationID string, params ...string) {
	// if you don't give a name use the name user have in his account
	if len(params) == 0 && user.Email != "" {
		params = strings.Split(user.Name, " ")
	}
	name := "there"
	snoozeDays := int64(3)

	if len(params) > 0 {
		val, err := strconv.Atoi(params[0])
		if err != nil {
			name = params[0]
		} else {
			name = strings.Split(user.Name, " ")[0]
			snoozeDays = int64(val)
		}
	}

	message := "Hi " + name + " 👋 \n \n <b>Our friends at Buildium support your Happy Inspector subscription. \n\n They can be reached at 888-414-1988, or by submitting a ticket through your Buildium account.</b>\n\nBuildium Support team are the best place to help you with this query as they understand your unique workflow and are trained in Happy Inspector 💫 \n \n Please also feel free to take a look through our FAQ on the Buildium integration:  \n https://intercom.help/happyco/frequently-asked-questions/buildium-integration-faq/faq-buildium-integration  \n Thanks!  \n HappyBot ☺"

	addReply(conversationID, message)
	assignConversation(conversationID, "1615207")
	snoozeConversation(conversationID, int64(3))
}
