package main

import (
	// "strings"
	"time"
	//"fmt"
)
func init() {
	commands["inspection"] = showRecentInspections
}
//********************************************Adding commands********************************************
func showRecentInspections(user User, conversationID string, params ...string) {

	inspectionsRec := getInspections(user.UserID)
	message := "<b>Showing recent inspections (max 5) created by "+ user.Name+" in last 30 days</b>"
	message += "\n"
	for _, inspection := range inspectionsRec {
		var url = "https://manage.happyco.com/folder/" + inspection.FolderID + "/inspections/" + inspection.ID
		var date, _ = time.Parse(time.RFC3339, inspection.CreatedAt)
		formattedDate := date.Format("02 Jan 2006 3:04PM")

		message += "<a href=" + url + ">" + url + "</a>" + " " + formattedDate
		message += "\n"
	}
	addNote(conversationID, message)
}
