package main
import (
	"encoding/json"
	//"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	// "strconv"
	"strings"
	// "time"

	//"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"

	// intercom "gopkg.in/intercom/intercom-go.v2"
)

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
	ID	string `db:"business_id"`
	Role string `db:"business_role_id"`
}
type IAP struct {
	Expiry string `db:"expires_at"`
}

// structs for reading payload in json received from Intercom

type ConversationMessage struct{
	Body string `json:"body"`
	Subject string `json:"subject"`
}
type User struct {
	UserID 	string `json:"user_id"`
	Type	string `json:"type"`
	Name string `json:"name"`
	Email string `json:"email"`
}

type Item struct {
	ConversationID 	string `json:"id"`
	User			User   `json:"user"`
	ConversationMessage ConversationMessage `json:"conversation_message"`
}

type Data struct {
	Item Item `json:"item"`
}

type Message struct {
	Data Data `json:"data"`
}
//********************************************New Conversation********************************************

//gets intercom token, admin list, reads the payload, and post note as a reply in the conversation
func newConversation(w http.ResponseWriter, r *http.Request) {
	fmt.Println("newConversation");
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
		user- user's attributes - name id email type
		conversationId - Intercom conversation ID
		conversationmessage- user's message
	*/

	user := msg.Data.Item.User
	conversationId := msg.Data.Item.ConversationID
	conversationMessage:=msg.Data.Item.ConversationMessage.Body
	conversationSubject:=msg.Data.Item.ConversationMessage.Subject

	go processNewConversation(user, conversationId, conversationMessage, conversationSubject)
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Received"))
}

func processNewConversation(user User, conversationID string, conversationMessage string, conversationSubject string) {
	fmt.Println("processNewConversation");
	// user.type - lead/user
	if user.Type == "user" {
		fmt.Println("message from a user. calling makeAndSendNote");
		makeAndSendNote(user, conversationID)

		// buildium autoresponder - only if not buildium support
		buildiumSupport := strings.Contains(user.Email,"@buildium.com")
		var ignorePhrases []string 
		var autoReply bool

		ignorePhrases:=["Automatic Reply","Automatic reply",,"automatic reply",,"Auto-reply",,"auto reply",
		,"Automatic Reply",,"Automatic Reply",,"Automatic Reply"]

		for _, phrase := range ignorePhrases{
			val:=strings.Contains(conversationSubject, phrase)

			if val {
				autoReply = true
				break
			}
		}

		planType := getUserPlanType(user.UserID)

		if planType == "buildium" && !buildiumSupport && !autoReply {
			sendBuildiumReply(user, conversationID)
		}
	}

	// change password autoresponder
	conversationMessage = strings.ToLower(conversationMessage)
	if strings.Contains(conversationMessage,"change password") {
		sendPasswordReply(user, conversationID)
	}
}

func getUserPlanType(ID string) (planType string) {
	fmt.Println("getUserPlanType");
	// DD/buildium/mri
	err := db.Get(&planType, "Select plan_type FROM current_subscriptions WHERE business_id IN (SELECT business_id from business_membership WHERE user_id = $1 AND inactivated_at IS NULL)", ID)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error in plan query %v: %v\n", ID, err)
	}
	return
}
