package lrreq

import dbreq "github.com/braginantonev/gcalc-server/pkg/database/requests-types"

const (
	CREATE_Table dbreq.DBRequestType = `
	CREATE TABLE IF NOT EXISTS users(
		id INTEGER PRIMARY KEY AUTOINCREMENT, 
		user TEXT,
		password TEXT
	);`

	// Required 2 args - username (string), encrypted password (string)
	INSERT_UserPass dbreq.DBRequestType = "INSERT INTO users (user, password) values ($1, $2)"

	// Required 1 arg - username (string)
	SELECT_UserPass dbreq.DBRequestType = "SELECT user, password FROM users WHERE user = $1"
)
