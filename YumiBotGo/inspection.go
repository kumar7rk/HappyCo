package main

import (
	"strconv"
	"time"
	"encoding/json"
)
type Asset struct {
	EmbeddedAsset EmbeddedAsset
}
type EmbeddedAsset struct {
	Address Address
	Building string
}
type Address struct {
	Line1 string 
	Line2 string 
	Locality string 
	Province string 
}

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
	message := "<b>Showing " + strconv.Itoa(len(inspectionsRec)) + " recent inspections created by " + user.Name + " in last 30 days</b>"
	message += "\n"
	for _, inspection := range inspectionsRec {
		var url = "https://manage.happyco.com/folder/" + inspection.FolderID + "/inspections/" + inspection.ID
		var date, _ = time.Parse(time.RFC3339, inspection.CreatedAt)
		var formattedDate = date.Format("02 Jan 2006 3:04PM")

		var msg Asset
    	var val []byte = []byte	(inspection.Asset)

	    _= json.Unmarshal(val, &msg)
	    
		message += inspection.FolderName +" > "
		var address string
		if msg.EmbeddedAsset.Address.Line1 != "" {
			address +=msg.EmbeddedAsset.Address.Line1+", "
		}
		if msg.EmbeddedAsset.Address.Line2 != ""{
			address +=msg.EmbeddedAsset.Address.Line2+", "
		}
		if msg.EmbeddedAsset.Address.Locality !="" {
			address +=msg.EmbeddedAsset.Address.Locality+", "
		}
		if msg.EmbeddedAsset.Address.Province != ""{
			address +=msg.EmbeddedAsset.Address.Province
		}

		message += "<a href=\"" + url + "\">" + address + "</a>" + " " + formattedDate
		message += "\n"

	}
	addNote(conversationID, message)
}
