package main

import (
	"strconv"
	"time"

	"happyco/libs/log"
)

func init() {
	getCommands["report"] = GetCommand{Func: showRecentReports, Description: `Shows reports created by a user in last 30 days. ` + "\n" + ` <b>Default</b> is 5. ` + "\n" + ` <b>Argument</b> # of reports to show. ` + "\n" + `<b>Example</b> yumi get report 100 will show 100 reports in last 30 days`}
}

//********************************************Adding commands********************************************
func showRecentReports(user User, conversationID string, params ...string) {
	var limit = 5

	if len(params) > 0 {
	var err error
		limit, err = strconv.Atoi(params[0])
		if err != nil {
			log.Error.KV("err", err).KV("params", params).KV("conversationID", conversationID).Println("could not parse number of snooze days for showing recent reports")
		}
	}
	reportsRec := getReports(user.UserID, limit)

	message := "<b>Showing " + strconv.Itoa(len(reportsRec)) + " recent reports created by " + user.Name + " in last 30 days</b>"
	message += "\n"
	for _, report := range reportsRec {
		var url = "https://manage.happyco.com/reports/" + report.PublicID
		var date, _ = time.Parse(time.RFC3339, report.CreatedAt)
		var formattedDate = date.Format("02 Jan 2006 3:04PM")

		message += "<a href=\"" + url + "\">" + report.Name + "</a>" + " " + formattedDate
		message += "\n"
	}
	addNote(conversationID, message)
}
