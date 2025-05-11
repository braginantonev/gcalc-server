package orchreq

import dbreq "github.com/braginantonev/gcalc-server/pkg/database/requests-types"

const (
	DBRequest_CREATE_Table dbreq.DBRequestType = `
	CREATE TABLE IF NOT EXISTS expressions(
		id INTEGER PRIMARY KEY AUTOINCREMENT, 
		user TEXT,
		internal_id INTEGER,
		string TEXT,
		status INTEGER,
		result REAL
	);`

	// Required 1 arg - user (string)
	DBRequest_SELECT_Expressions dbreq.DBRequestType = "SELECT user, internal_id, string, status, result FROM expressions WHERE user = $1"

	// Required 5 args - user (string), internal_id (int32), str (string), status (int32), result (float64)
	DBRequest_INSERT_Expression dbreq.DBRequestType = "INSERT INTO expressions (user, internal_id, string, status, result) values ($1, $2, $3, $4, $5)"

	// Required 2 args - user (string), internal_id (int32)
	DBRequest_SELECT_Expression dbreq.DBRequestType = "SELECT user, internal_id, string, status, result FROM expressions WHERE user = $1 AND internal_id = $2"

	// Required 4 args - status (int32), result (float64), user (string), internal_id (int32)
	DBRequest_UPDATE_Expression dbreq.DBRequestType = "UPDATE expressions SET status = $1, result = $2 WHERE user = $3 AND internal_id = $4"
)
