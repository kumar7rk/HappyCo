package main

import (
	"encoding/json"
	"strconv"
	"time"

	"happyco/libs/log"
)

type Asset struct {
	EmbeddedAsset EmbeddedAsset
}
type EmbeddedAsset struct {
	Address  Address
	Building string
}
type Address struct {
	Line1    string
	Line2    string
	Locality string
	Province string
}

//********************************************Init********************************************
func init() {
	getCommands["inspection"] = GetCommand{Func: showRecentInspections, Description: `Shows inspections created by a user in last 30 days. ` + "\n" + ` <b>Default</b> is 5. ` + "\n" + ` <b>Argument</b> # of inspections to show. ` + "\n" + `<b>Example</b> yumi get inspection 100 will show 100 inspections in last 30 days`}
}

//********************************************Adding commands********************************************
func showRecentInspections(user User, conversationID string, params ...string) {
	var limit = 5
	var err error
	
	if len(params) > 0 {
		limit, err = strconv.Atoi(params[0])
		if err != nil {
			log.Error.KV("err",err).KV("params", params).KV("conversationID", conversationID).Println("could not parse number of snooze days for showing recent inspections")
		}
	}
	inspectionsRec := getInspections(user.UserID, limit)
	message := "<b>Showing " + strconv.Itoa(len(inspectionsRec)) + " recent inspections created by " + user.Name + " in last 30 days</b>"
	message += "\n"
	for _, inspection := range inspectionsRec {
		var url = "https://manage.happyco.com/folder/" + inspection.FolderID + "/inspections/" + inspection.ID
		var date, _ = time.Parse(time.RFC3339, inspection.CreatedAt)
		var formattedDate = date.Format("02 Jan 2006 3:04PM")

		var msg Asset
		var val []byte = []byte(inspection.Asset)

		_ = json.Unmarshal(val, &msg)

		// message += "<b>" + inspection.FolderName +"</b> "
		var address string
		if msg.EmbeddedAsset.Address.Line1 != "" {
			address += msg.EmbeddedAsset.Address.Line1 + ", "
		}
		if msg.EmbeddedAsset.Address.Line2 != "" {
			address += msg.EmbeddedAsset.Address.Line2 + ", "
		}
		if msg.EmbeddedAsset.Address.Locality != "" {
			address += msg.EmbeddedAsset.Address.Locality + ", "
		}
		if msg.EmbeddedAsset.Address.Province != "" {
			address += msg.EmbeddedAsset.Address.Province
		}
		// message += "<p><a href=\"" + url + "\"><p>" + inspection.FolderName + "</p><p>" + address + "</p><p>" + formattedDate + "</p></a>\n"
		message += inspection.FolderName + " > "
		message += "<a href=\"" + url + "\">" + address + "</a>" + " " + formattedDate + "\n\n"
	}
	addNote(conversationID, message)
}
