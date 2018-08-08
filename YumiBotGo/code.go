/* Adding User Data from DB to Intercom as a note for every new conversation */

package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"

	intercom "gopkg.in/intercom/intercom-go.v2"

	"happyco/apps/tools/libs/buildvars"
)

//********************************************Variable declaration********************************************

type Inspection struct {
	Business     string
	User         string
	Role         string
	FolderID     string `db:"folder_id"`
	FolderName   string `db:"folder_name"`
	CreatedAt    string `db:"created_at"`
	TemplateName string `db:"template_name"`
	ID           string
	Status       string
	Location     string
}

type Report struct {
	Business   string
	User       string
	Role       string
	FolderID   string `db:"folder_id"`
	FolderName string `db:"folder_name"`
	CreatedAt  string `db:"created_at"`
	Name       string
	PublicID   string `db:"public_id"`
	Location   string
}

type Business struct {
	ID   string `db:"business_id"`
	Role string `db:"business_role_id"`
}
type IAP struct {
	Expiry string `db:"expires_at"`
}
type Plan struct {
	Type string `db:"plan_type"`
}
var db *sqlx.DB

// structs for reading payload in json received from Intercom
type User struct {
	UserID string `json:"user_id"`
	Type   string `json:"type"`
}

type Item struct {
	ConversationID string `json:"id"`
	User           User   `json:"user"`
}

type Data struct {
	Item Item `json:"item"`
}

type Message struct {
	Data Data `json:"data"`
}

//********************************************Main function********************************************

func main() {
	var err error

	postgresURI := os.Getenv("POSTGRES_URI")
	if postgresURI == "" {
		fmt.Println("URI error")
	}

	db, err = sqlx.Connect("postgres", postgresURI)

	if err == nil {
		fmt.Println("DB Connected")
	}
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error connecting to DB%v: %v\n", "", err)
		return
	}

	//handling every new convesations in newConversation method
	http.HandleFunc("/", newConversation)
	http.HandleFunc("/healthcheck", healthcheck)
	if err := http.ListenAndServe(":8080", nil); err != nil {
		panic(err)
	}
}

//********************************************New Conversation********************************************

//gets intercom token, admin list, reads the payload, and post note as a reply in the conversation
func newConversation(w http.ResponseWriter, r *http.Request) {

	// Read body/payload
	b, err := ioutil.ReadAll(r.Body)
	defer r.Body.Close()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Unmarshal the json
	var msg Message
	err = json.Unmarshal(b, &msg)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	/* getting attributes from the received json
	user type - lead/user
	userId - happyCo user id
	conversationId - Intercom conversation Id
	*/

	userType := msg.Data.Item.User.Type
	userId := msg.Data.Item.User.UserID
	conversationId := msg.Data.Item.ConversationID

	//only run the following code when the received message is from a HappyCo user
	if userType == "user" {
		go makeAndSendNote(userId, conversationId)
	}
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Received"))
}

func makeAndSendNote(ID string, conversationID string) {
	p := fmt.Println

	// gets intercom access token
	accessToken := os.Getenv("INTERCOM_ACCESS_TOKEN")
	ic := intercom.NewClient(accessToken, "")

	user, err := ic.Users.FindByUserID(ID)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error while finding user ID %v: %v\n", ID, err)
		return
	}

	//testing prints
	p("Conversation id: " + conversationID)
	p("User id: " + ID)
	p("User name: " + user.Name)

	// calling the method to compile the note with all the required information
	note, plan_type := makeNote(ID)
	_, err = ic.Conversations.Reply(conversationID, intercom.Admin{ID: "207278"}, intercom.CONVERSATION_NOTE, note)
	//copied and pasted from api-docs
	if herr, ok := err.(intercom.IntercomError); ok && herr.GetCode() == "not_found" {
		fmt.Fprintf(os.Stderr, "Error from Intercom when adding note %v: %v\n", "", err)
		return
	}

	if plan_type == "buildium" {

		//extracting firstName from the user name.
		firstName := strings.Fields(user.Name)

		buildiumMessage := "Hi " + firstName[0] + " 👋 \n \n <b>Our friends at Buildium support your Happy Inspector subscription. \n\n They can be reached at 888-414-1988, or by submitting a ticket through your Buildium account.</b>\n\nBuildium Support team are the best place to help you with this query as they understand your unique workflow and are trained in Happy Inspector 💫 \n Please also feel free to take a look through our FAQ on the Buildium integration:  \n https://intercom.help/happyco/frequently-asked-questions/buildium-integration-faq/faq-buildium-integration  \n Thanks!  \n HappyCo team ☺"
		_, err = ic.Conversations.Reply(conversationID, intercom.Admin{ID: "207278"}, intercom.CONVERSATION_COMMENT, buildiumMessage)

		if herr, ok := err.(intercom.IntercomError); ok && herr.GetCode() == "not_found" {
			fmt.Fprintf(os.Stderr, "Error from Intercom when replying to Buildium %v: %v\n", "", err)
			return
		}
	}
}

//********************************************Getting UserData********************************************

//queries the db and adds returned values in array
func getUserData(ID string) (inspectionsRec []Inspection, reportsRec []Report, businessRec []Business, iapRec []IAP, integrationName string, planTypeRec []Plan) {

	//fetching most recent (5) inspections for the user within the last 30 days.
	err := db.Select(&inspectionsRec, "SELECT folders.business,folders.user,folders.role ,folders.folder_id,folders.folder_name,i.created_at as created_at,i.template_name,i.id,i.status,i.location FROM (SELECT businesses.business_id as business,businesses.user_id as user,role_id as role,folder_id as folder_id,folder_name as folder_name FROM (SELECT bm.business_id as business_id,bm.user_id as user_id,bm.business_role_id as role_id,f.id as folder_id,f.name as folder_name FROM business_membership as bm JOIN portfolios as f ON bm.business_id = f.business_id WHERE bm.user_id = $1 AND bm.inactivated_at IS NULL AND f.inactivated_at IS NULL) as businesses GROUP BY businesses.business_id,businesses.role_id,businesses.user_id,folder_id,folder_name ORDER BY businesses.business_id ) as folders JOIN inspections as i ON folders.folder_id = i.folder_id WHERE i.user_id = $1::varchar AND i.archived_at IS NULL AND i.created_at > (CURRENT_DATE- interval '30 day') ORDER BY i.created_at DESC LIMIT $2", ID, 5)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error in inspection query %v: %v\n", ID, err)
	}

	//fetching most recent (5) reports for the user within the last 30 days.
	err = db.Select(&reportsRec, "SELECT folders.business,folders.user,folders.role,folders.folder_id,folders.folder_name,r.created_at as created_at,r.name,r.public_id,r.location FROM (SELECT businesses.business_id as business,businesses.user_id as user,role_id as role,folder_id as folder_id,folder_name as folder_name FROM (SELECT bm.business_id as business_id,bm.user_id as user_id,bm.business_role_id as role_id,f.id as folder_id,f.name as folder_name FROM business_membership as bm JOIN portfolios as f ON bm.business_id = f.business_id WHERE bm.user_id = $1 AND bm.inactivated_at IS NULL AND f.inactivated_at IS NULL) as businesses GROUP BY businesses.business_id,businesses.role_id,businesses.user_id,folder_id,folder_name ORDER BY businesses.business_id ) as folders JOIN reports_v3 as r ON folders.folder_id = r.folder_id WHERE r.user_id = $1::varchar AND r.archived_at IS NULL AND r.created_at > (CURRENT_DATE- interval '30 day') ORDER BY r.created_at DESC LIMIT $2", ID, 5)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error in reoprt query %v: %v\n", ID, err)
	}
	// fetching business id and role id for user role in this business
	err = db.Select(&businessRec, "SELECT business_id,business_role_id FROM business_membership WHERE user_id = $1 AND inactivated_at IS NULL", ID)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error in business query %v: %v\n", ID, err)
	}
	// checking if the business is on IAP get the expiry date
	err = db.Select(&iapRec, "SELECT expires_at FROM iap_receipts WHERE company_id IN (SELECT business_id FROM business_membership WHERE user_id = $1) ORDER BY expires_at DESC limit 1", ID)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error in iap query %v: %v\n", ID, err)
	}

	// Check if the business has integration w/Yardi
	var integrationCount int

	err = db.Get(&integrationCount, "Select COUNT(*) FROM integration_yardi_properties WHERE business_id IN (SELECT business_id from business_membership WHERE user_id = $1 AND inactivated_at IS NULL)", ID)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error in integration query-Yardi %v: %v\n", ID, err)
	}
	if integrationCount > 0 {
		integrationName = "Yardi"
	}
	// MRI
	err = db.Get(&integrationCount, "Select COUNT(*) FROM integration_mri_properties WHERE business_id IN (SELECT business_id from business_membership WHERE user_id = $1 AND inactivated_at IS NULL)", ID)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error in integration query-MRI %v: %v\n", ID, err)
	}

	if integrationCount > 0 {
		integrationName = "MRI"
	}
	// Resman
	err = db.Get(&integrationCount, "Select COUNT(*) FROM integration_resman_properties WHERE business_id IN (SELECT business_id from business_membership WHERE user_id = $1 AND inactivated_at IS NULL)", ID)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error in integration query-Resman %v: %v\n", ID, err)
	}
	if integrationCount > 0 {
		integrationName = "Resman"
	}
	var b_id = "36439"

	// DD/buildium/mri
	err = db.Select(&planTypeRec, "Select plan_type FROM current_subscriptions WHERE business_id IN (SELECT business_id from business_membership WHERE user_id = $1 AND inactivated_at IS NULL)", b_id)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error in plan query %v: %v\n", ID, err)
	}
	return
}

//********************************************Building note********************************************

// code starts running from here.
// build the note in a string format
// should be called when a new intercom message is received
func makeNote(us_id string) (string, string) {
	var note string
	var formattedDate string

	//getting user data from the database
	inspectionsRec, reportsRec, businessRec, iapRec, integrationName, planTypeRec := getUserData(us_id)

	//******************constructing business string******************
	note = "<b>A small note from Yumi 🐶</b><br/><br/>"
	note+="\n\n Test Note"

	note += "<b>✅User is associated with the following businesses</b><br/><br/>"

	// 1, 2, 3, 4 = Constant Admin, PM, Inspector, Limited Inspector
	// 8, 9 = Basic Admin, PM
	var roles = []string{"", "Admin", "Process Manger", "Inspector", "Limited Inspector", "", "", "", "Admin", "Process Manager"}

	for _, business := range businessRec {

		permission, err := strconv.Atoi(business.Role)
		if err != nil {
			panic(err)
		}
		var roleID = business.Role

		if roleID == "1" || roleID == "2" || roleID == "3" || roleID == "4" {
			note += " The business is on Constant/full Permissions \n"
		}

		if roleID == "8" || roleID == "9" {
			note += " The business is on Basic Permissions \n"
		}

		note += roles[permission] + " for "
		var text = "https://manage.happyco.com/admin/businesses/" + business.ID
		note += text + "<br/><br/>"
	}
	note += "\n"

	//******************constructing integration string******************
	if integrationName != "" {
		note += "\n"
		note += "\n"
		note += "The business is " + integrationName
	}

	//******************constructing plan type string******************
	planType :="plan type"
	for _, plan := range planTypeRec {
		if plan.Type == "due_diligence" {
			note += "\n"
			note += "\n"
			note += "Plan: " + "Due Diligence"
		}
		if plan.Type == "buildium" {
			note += "\n"
			note += "\n"
			note += "Plan: " + "Buildium"
			planType = plan.Type
		}
		if plan.Type == "mri" {
			note += "\n"
			note += "\n"
			note += "Plan: " + "MRI"
		}
	}
	note += "\n"
	note += "\n"

	//******************constructing inspection string******************
	note += "<b>✅   Yumi found these recent (max: 5) <em>Inspections in last 30 days:</em></b><br/>"
	note += "\n"
	var url string
	for _, inspection := range inspectionsRec {

		url = "https://manage.happyco.com/folder/" + inspection.FolderID + "/inspections/" + inspection.ID

		var date, _ = time.Parse(time.RFC3339, inspection.CreatedAt)
		formattedDate = date.Format("02 Jan 2006 3:04PM")

		note += "<a href=" + url + ">" + url + "</a>" + " " + formattedDate
		note += "\n"
	}
	//******************constructing report string******************

	note += "\n"
	note += "\n"
	note += "<b>✅   Yumi found these recent (max: 5) <em>Reports in last 30 days:</em></b><br/>"
	note += "\n"
	for _, report := range reportsRec {

		var url = "https://manage.happyco.com/reports/" + report.PublicID

		var date, _ = time.Parse(time.RFC3339, report.CreatedAt)
		formattedDate = date.Format("02 Jan 2006 3:04PM")

		note += "<a href=" + url + ">" + report.Name + "</a>" + " " + formattedDate
		note += "\n"
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
	note +="\n"
	note +="\n"

	note += "<a href=" + "https://hpy.io/yumi" + ">" + "Feedback/Report issue" + "</a>"
	return note, planType
}

func healthcheck(w http.ResponseWriter, r *http.Request) {
	if ping() != nil {
		w.WriteHeader(http.StatusServiceUnavailable)
		w.Write([]byte("Database timeout"))
		return
	}
	w.Write([]byte("OK " + buildvars.BuildScmRevisionShort + " " + buildvars.BuildScmStatus))
}

func ping() error {
	ping := make(chan error, 0)
	timeout := time.After(10 * time.Second)
	go func() {
		ping <- db.Ping()
	}()
	var err error
	select {
	case <-timeout:
		err = errors.New("Postgres ping timeout error")
	case err = <-ping:
	}
	return err
}
