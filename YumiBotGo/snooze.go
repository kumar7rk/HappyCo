package main

import (
	"bytes"
	"io/ioutil"
	"net/http"
	"strconv"
	"time"
    "os"
)

func snoozeConversation(conversationID string, days int64) {
	currentTimeInSec := int64(time.Now().Unix())
	secInADay := int64(86400)
	snoozeTimeInSec := int64(secInADay*days)

	snooze_until := currentTimeInSec + snoozeTimeInSec

	url := "https://api.intercom.io/conversations/" + conversationID + "/reply"
	payload := []byte(`{ "admin_id":"207278", "message_type":"snoozed", "snoozed_until":` + strconv.FormatInt(snooze_until, 10) + `}`)

	req, _ := http.NewRequest("POST", url, bytes.NewBuffer(payload))
	req.Header.Set("Authorization", "Bearer "+ os.Getenv("INTERCOM_ACCESS_TOKEN"))
	req.Header.Set("Accept", "application/json")
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, _ := client.Do(req)

    //MT- not sure about error handling here 
	/*if err != nil {
        fmt.Fprintf(os.Stderr, "Error in Snoozing conversation %v: %v\n", ID, err)
    }*/
    
    defer resp.Body.Close()
	
    _, _ = ioutil.ReadAll(resp.Body)
}
