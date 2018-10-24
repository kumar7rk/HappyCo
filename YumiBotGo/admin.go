package main

import (
	"strings"
)

//********************************************Init********************************************
func init() {
	getCommands["admin"] = GetCommand{Func: showAllAdmins, Description: `Prints a note with the name and emails of all the admins in a business.
		<b>Arguments</b> userID.
		<b>Example</b> yumi get admin [userID]`}
}

//********************************************Adding commands********************************************
func showAllAdmins(user User, conversationID string, params ...string) {
	var message string
	userID := user.UserID
	if len(params) > 0 {
		userID = params[0]
	}
	adminsRec := getAdmins(userID)
	for _, admin := range adminsRec {
		if strings.Contains(admin.Detail, "hpy.io") || strings.Contains(admin.Detail, "happy.co") || strings.Contains(admin.Detail, "happyco.com") {
			continue
		}
		message += admin.Detail
		message += "\n"
	}
	if message == "" {
		message += "No non-HappyCo admin found"
	}
	addNote(conversationID, message)
}
