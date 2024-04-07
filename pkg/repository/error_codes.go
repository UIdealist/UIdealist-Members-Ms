package repository

const (

	// Error codes
	DATABASE_ERROR             string = "connection_database_error"
	CACHE_ERROR                string = "connection_cache_error"
	USER_NOT_FOUND             string = "auth_user_not_found"
	INVALID_CREDENTIALS        string = "auth_invalid_credentials"
	INVALID_DATA               string = "auth_invalid_form_data"
	TEAM_NOT_FOUND             string = "team_not_found"
	TEAM_MEMBER_NOT_FOUND      string = "team_member_not_found"
	TEAM_MEMBER_ALREADY_EXISTS string = "team_member_already_exists"
	TEAM_MEMBER_NOT_REMOVED    string = "team_member_not_removed"
	TEAM_MEMBER_NOT_ADDED      string = "team_member_not_added"
	TEAM_MEMBER_NOT_LISTED     string = "team_member_not_listed"
	TEAM_MEMBER_NOT_UPDATED    string = "team_member_not_updated"
	NOT_AUTHORIZED             string = "auth_not_authorized"
)
