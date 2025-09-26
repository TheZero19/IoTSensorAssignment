package Constants

import "database/sql"

var Db *sql.DB

const AUTH_HEADER_KEY string = "Authorization"
const AUTH_HEADER_VALUE_SEPARATOR string = " "
const API_KEY_AUTH_REGISTRATION_TYPE_PREFIX string = "API-KEY"
const PSK_AUTH_TYPE_PREFIX string = "PSK"
