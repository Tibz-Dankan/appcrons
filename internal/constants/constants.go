package constants

var QUERY_MINIMUM_LIMIT float64 = 5
var QUERY_MAXIMUM_LIMIT float64 = 30

var RECORD_NOT_FOUND_ERROR = "record not found"
var NO_FILE_UPLOADED_ERROR = "there is no uploaded file associated with the given key"

// UNKNOWN_USER_EMAIL is the sentinel email for the shared placeholder
// "unknown user" row that anonymous-tolerant routes (e.g. site visits)
// attribute their userId to instead of leaving it blank.
var UNKNOWN_USER_EMAIL = "unknown.user@appcrons.internal"
