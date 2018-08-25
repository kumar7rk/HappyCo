package main

import (
	"fmt"
	"os"
)

func getInspections(ID string) (inspectionsRec []Inspection) {
	err := db.Select(&inspectionsRec, "SELECT folders.business,folders.user,folders.role ,folders.folder_id,folders.folder_name,i.created_at as created_at,i.template_name,i.id,i.status,i.location FROM (SELECT businesses.business_id as business,businesses.user_id as user,role_id as role,folder_id as folder_id,folder_name as folder_name FROM (SELECT bm.business_id as business_id,bm.user_id as user_id,bm.business_role_id as role_id,f.id as folder_id,f.name as folder_name FROM business_membership as bm JOIN portfolios as f ON bm.business_id = f.business_id WHERE bm.user_id = $1 AND bm.inactivated_at IS NULL AND f.inactivated_at IS NULL) as businesses GROUP BY businesses.business_id,businesses.role_id,businesses.user_id,folder_id,folder_name ORDER BY businesses.business_id ) as folders JOIN inspections as i ON folders.folder_id = i.folder_id WHERE i.user_id = $1::varchar AND i.archived_at IS NULL AND i.created_at > (CURRENT_DATE- interval '30 day') ORDER BY i.created_at DESC LIMIT $2", ID, 5)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error in inspection query %v: %v\n", ID, err)
	}
	return
}