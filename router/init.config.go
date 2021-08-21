package router

const (
	IS_RELEASE                 bool = false
	SERVER_PORT                string = "55699"
	SWAGGER_API_DEFINITION_URL string = "http://localhost:55699/swagger/doc.json" // The url pointing to API definition
	TOKEN_VALID_TIME_IN_SECOND int64  = 10 * 24 * 60 * 60                         // 10 days
)
