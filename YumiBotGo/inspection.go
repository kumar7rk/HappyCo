package main

import (
	"strconv"
	"time"
)
//********************************************Init********************************************
func init() {
	commands["inspection"] = showRecentInspections
}

//********************************************Adding commands********************************************
func showRecentInspections(user User, conversationID string, params ...string) {
	var limit = 5
	if len(params) > 0 {
		limit, _ = strconv.Atoi(params[0])
	}
	inspectionsRec := getInspections(user.UserID, limit)
	message := "<b>Showing " + strconv.Itoa(limit) + " recent inspections created by " + user.Name + " in last 30 days</b>"
	message += "\n"
	for _, inspection := range inspectionsRec {
		var url = "https://manage.happyco.com/folder/" + inspection.FolderID + "/inspections/" + inspection.ID
		var date, _ = time.Parse(time.RFC3339, inspection.CreatedAt)
		var formattedDate = date.Format("02 Jan 2006 3:04PM")

		message += "<a href=\"" + url + "\">" + url + "</a>" + " " + formattedDate
		message += "\n"
	}
	addNote(conversationID, message)
}
