package model

const (
	DB_USER     string = "root"
	DB_PASSWORD string = "root"
	DB_HOST     string = "localhost"
	DB_PORT     string = "3306"
	DB_NAME     string = "cmsdiy"
	DB_CHARSET  string = "utf8mb4"

	ROOT_ROLENAME string = "author"
	ROOT_USERNAME string = "jinyaoMa"
	ROOT_ACCOUNT  string = "root"
	ROOT_PASSWORD string = "root"

	ROLE_DEFAULT_MEMBER_NAME       string = "member"
	ROLE_DEFAULT_MEMBER_PERMISSION string = ""
	ROLE_DEFAULT_MEMBER_SPACE      Size   = 0
	ROLE_DEFAULT_MEMBER_CODE       string = "123456789012"

	ROLE_CODE_SIZE  int = 12
	SHARE_CODE_SIZE int = 4
	JWT_KEY_SIZE    int = 16
)

var (
	StorageBranches []string = []string{"D:/_cmsdiy0", "D:/_cmsdiy1"}
	UserLimit       int64    = 5
	WorkspaceLimit  Size     = 280 * GB
)
