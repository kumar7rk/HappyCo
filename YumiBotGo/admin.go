package main

import (
	"fmt"
)

//********************************************Init********************************************
func init() {
	commands["admin"] = showAllAdmins
}

//********************************************Adding commands********************************************
func showAllAdmins(user User, conversationID string, params ...string) {
	var message string
	fmt.Println("showAllAdmins")
	adminsRec := getAdmins(user.UserID)
	for _, admin := range adminsRec {
		message += admin.Detail
		message +="\n"
	}
	addNote(conversationID, message)
}
