package models

import "bank-api/db"

func OperationTypeExists(id int) bool {
	var exists bool
	err := db.DB.QueryRow(`SELECT EXISTS(SELECT 1 FROM operation_types WHERE operation_type_id = $1)`, id).Scan(&exists)
	if err != nil {
		return false
	}
	return exists
}
