package models

const (
	SERVER_PORT string = "55699"

	DB_USER     string = "root"
	DB_PASSWORD string = "root"
	DB_HOST     string = "127.0.0.1"
	DB_PORT     string = "7531"
	DB_NAME     string = "cmsdiy"
	DB_CHARSET  string = "utf8mb4"

	ROOT_ROLENAME string = "author"
	ROOT_USERNAME string = "jinyaoMa"
	ROOT_ACCOUNT  string = "root"
	ROOT_PASSWORD string = "root"

	ROLE_CODE_SIZE  int = 16
	SHARE_CODE_SIZE int = 4
	JWT_KEY_SIZE    int = 12
)

var (
	StorageBranches []string = []string{"D:/WORKSPACE", "D:/TESTSPACE"}
	UserLimit       uint     = 5
	WorkspaceLimit  Size     = 280 * GB
)
