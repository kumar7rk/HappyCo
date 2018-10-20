package main

import (
	"fmt"
	"strconv"
	"strings"
	"time"

	_ "github.com/lib/pq"
)

//********************************************Init********************************************
func init() {
	getCommands["note"] = Command{Func: makeAndSendNote, Description: `Print user's information - business, role etc.`}
}

//********************************************Adding note to conversation********************************************
func makeAndSendNote(user User, conversationID string, params ...string) {
	ID := user.UserID
	if len(params) > 0 {
		ID = params[0]
	}

	p := fmt.Println

	//testing prints
	p("Conversation id: " + conversationID)
	p("User id: " + ID)
	p("User name: " + user.Name)
	p("User email: " + user.Email)

	note := makeNote(ID)

	addNote(conversationID, note)
}

//********************************************Getting UserData********************************************
func getUserData(ID string) (businessRec []Business, iapRec []IAP, integrationName string, planRecs []Plan) {
	// fetching business id and role id for user role in this business
	businessRec = getBusiness(ID)

	// checking if the business is on IAP get the expiry date
	iapRec = getIAP(ID)

	// Check if the business has integration w/Yardi
	integrationName = getIntegration(ID)

	// Getting Plan type - Admin tags
	planRecs = getUserPlans(ID)
	return
}

//********************************************Building note********************************************
func makeNote(us_id string) string {
	var note string
	var formattedDate string

	//getting user data from the database
	businessRec, iapRec, integrationName, planRecs := getUserData(us_id)

	note = "<b>🐶Note</b><br/><br/>"

	//******************constructing business string******************
	// 1, 2, 3, 4 = Constant Admin, PM, Inspector, Limited Inspector
	// 8, 9 = Basic Admin, PM
	var roles = []string{"", "Admin", "Process Manger", "Inspector", "Limited Inspector", "", "", "", "Admin", "Process Manager"}

	for _, business := range businessRec {
		userRoleID := business.Role.String
		var businessPermission = business.PermissionsModel

		if !business.Role.Valid {
			if businessPermission == "constant-roles" {
				userRoleID = "4"
			}
			if businessPermission == "basic-roles" {
				userRoleID = "9"
			}
		}
		permission, err := strconv.Atoi(userRoleID)
		if err != nil {
			fmt.Printf("Error converting user role id: %v\n", err)
		}

		if businessPermission == "constant-roles" {
			businessPermission = "Constant/Full"
		}

		if businessPermission == "basic-roles" {
			businessPermission = "Basic"
		}
		MRR := business.MRR.String
		if !business.MRR.Valid {
			MRR = "NA"
		}

		supportLevel := business.SupportLevel.String
		if !business.SupportLevel.Valid {
			supportLevel = "NA"
		}
		note += "<b>Business: </b>" + business.Name + "\n"
		note += "<b>BusinessID:</b>" + business.ID + "\n"
		note += "<b>Permissions:</b>" + businessPermission + "\n"
		note += "<b>Role:</b>" + roles[permission] + "\n"
		note += "<b>MRR:</b>" + MRR + "\n"
		if !business.SupportLevel.Valid {
			note += "<b>Support Level:</b>" + supportLevel + "\n"
		} else {
			note += "<h2><b>Support Level:</b>" + supportLevel + "</h2>\n"
		}
	}
	//******************constructing plan type string******************
	for _, plan := range planRecs {
		note += "<b>Plan: </b>" + strings.Title(plan.Name) + "\n"

		if plan.Status != "active" {
			note += "<b><h3>Status:</b>" + plan.Status + "</h3>\n"
		}
	}
	//******************constructing integration string******************
	if integrationName != "" {
		note += "<b>Integration: </b>" + integrationName
	}

	//******************constructing iap string******************
	for _, iap := range iapRec {
		if iap.Expiry != "" {
			var date, _ = time.Parse(time.RFC3339, iap.Expiry)
			formattedDate = date.Format("02 Jan 2006 3:04PM")
			note += "\n"
			note += "\n"
			note += "<b>The business is on IAP. It expires on </b>" + formattedDate
		}
	}
	note += "\n"
	note += "<a href=" + "https://hpy.io/yumi" + ">" + "Feedback/Report incorrect information" + "</a>"
	return note
}
