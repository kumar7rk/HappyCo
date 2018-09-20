package main

import (
	"fmt"
	"os"
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
	Asset        string
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
	ID               string 
	Name             string 
	Role             string `db:"business_role_id"`
	PermissionsModel string `db:"permissions_model"`
}
type IAP struct {
	Expiry string `db:"expires_at"`
}
type Plan struct {
	Type   string `db:"plan_type"`
	Status string 
}
type Admin struct {
	Detail string
}

//********************************************Inspection********************************************

func getInspections(userID string, limit int) (inspectionsRec []Inspection) {
	err := db.Select(&inspectionsRec, "SELECT folders.business,folders.user,folders.role ,folders.folder_id,folders.folder_name,i.created_at as created_at,i.template_name,i.id,i.status,i.location,i.asset FROM (SELECT businesses.business_id as business,businesses.user_id as user,role_id as role,folder_id as folder_id,folder_name as folder_name FROM (SELECT bm.business_id as business_id,bm.user_id as user_id,bm.business_role_id as role_id,f.id as folder_id,f.name as folder_name FROM business_membership as bm JOIN portfolios as f ON bm.business_id = f.business_id WHERE bm.user_id = $1 AND bm.inactivated_at IS NULL AND f.inactivated_at IS NULL) as businesses GROUP BY businesses.business_id,businesses.role_id,businesses.user_id,folder_id,folder_name ORDER BY businesses.business_id ) as folders JOIN inspections as i ON folders.folder_id = i.folder_id WHERE i.user_id = $1::varchar AND i.archived_at IS NULL AND i.created_at > (CURRENT_DATE- interval '30 day') ORDER BY i.created_at DESC LIMIT $2", userID, limit)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error in inspection query %v: %v\n", userID, err)
	}
	return
}

//********************************************Report********************************************
func getReports(userID string, limit int) (reportsRec []Report) {
	err := db.Select(&reportsRec, "SELECT folders.business,folders.user,folders.role,folders.folder_id,folders.folder_name,r.created_at as created_at,r.name,r.public_id,r.location FROM (SELECT businesses.business_id as business,businesses.user_id as user,role_id as role,folder_id as folder_id,folder_name as folder_name FROM (SELECT bm.business_id as business_id,bm.user_id as user_id,bm.business_role_id as role_id,f.id as folder_id,f.name as folder_name FROM business_membership as bm JOIN portfolios as f ON bm.business_id = f.business_id WHERE bm.user_id = $1 AND bm.inactivated_at IS NULL AND f.inactivated_at IS NULL) as businesses GROUP BY businesses.business_id,businesses.role_id,businesses.user_id,folder_id,folder_name ORDER BY businesses.business_id ) as folders JOIN reports_v3 as r ON folders.folder_id = r.folder_id WHERE r.user_id = $1::varchar AND r.archived_at IS NULL AND r.created_at > (CURRENT_DATE- interval '30 day') ORDER BY r.created_at DESC LIMIT $2", userID, limit)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error in reoprt query %v: %v\n", userID, err)
	}
	return
}

//********************************************Business********************************************
func getBusiness(userID string) (businessRec []Business) {

	err := db.Select(&businessRec, "SELECT b.id, b.name, bm.business_role_id, bc.permissions_model FROM businesses b INNER JOIN business_membership bm ON b.id = bm.business_id INNER JOIN business_customizations bc ON bc.business_id = bm.business_id WHERE bm.user_id = $1 AND inactivated_at IS NULL", userID)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error in business query %v: %v\n", userID, err)
	}
	return
}

//********************************************IAP********************************************
func getIAP(userID string) (iapRec []IAP) {
	err := db.Select(&iapRec, "SELECT expires_at FROM iap_receipts WHERE company_id IN (SELECT business_id FROM business_membership WHERE user_id = $1) ORDER BY expires_at DESC limit 1", userID)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error in iap query %v: %v\n", userID, err)
	}
	return
}

//********************************************Integration********************************************
func getIntegration(userID string) (integrationName string) {
	var integrationCount int
	err := db.Get(&integrationCount, "Select COUNT(*) FROM integration_yardi_properties WHERE business_id IN (SELECT business_id from business_membership WHERE user_id = $1 AND inactivated_at IS NULL)", userID)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error in integration query-Yardi %v: %v\n", userID, err)
	}
	if integrationCount > 0 {
		integrationName = "Yardi"
	}
	// MRI
	err = db.Get(&integrationCount, "Select COUNT(*) FROM integration_mri_properties WHERE business_id IN (SELECT business_id from business_membership WHERE user_id = $1 AND inactivated_at IS NULL)", userID)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error in integration query-MRI %v: %v\n", userID, err)
	}

	if integrationCount > 0 {
		integrationName = "MRI"
	}
	// Resman
	err = db.Get(&integrationCount, "Select COUNT(*) FROM integration_resman_properties WHERE business_id IN (SELECT business_id from business_membership WHERE user_id = $1 AND inactivated_at IS NULL)", userID)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error in integration query-Resman %v: %v\n", userID, err)
	}
	if integrationCount > 0 {
		integrationName = "Resman"
	}
	return
}

//********************************************Plan Type********************************************
func getUserPlanType(userID string) (planTypeRec []Plan) {
	err := db.Select(&planTypeRec, "Select plan_type,status FROM current_subscriptions WHERE business_id IN (SELECT business_id from business_membership WHERE user_id = $1 AND inactivated_at IS NULL)", userID)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error in plan query %v: %v\n", userID, err)
	}
	return
}

//********************************************Business's Admins********************************************
func getAdmins(userID string) (AdminRec []Admin) {
	err := db.Select(&AdminRec, "SELECT CONCAT(first_name, ' ', last_name, ' ', email) as detail FROM users u INNER JOIN business_membership bm ON u.id = bm.user_id AND bm.inactivated_at IS NULL AND bm.business_role_id IN (8, 1) INNER JOIN business_membership bm2 ON bm2.business_id = bm.business_id WHERE bm2.user_id =$1", userID)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error in plan query %v: %v\n", userID, err)
	}
	return
}
