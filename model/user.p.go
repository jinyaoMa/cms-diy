package model

func prepareRootUser() {
	_, hasRoot := GetUserByAccount(ROOT_ACCOUNT)
	role, hasRole := GetRoleByName(ROOT_ROLENAME)
	if !hasRoot && hasRole {
		user := User{
			Name:       ROOT_USERNAME,
			Account:    ROOT_ACCOUNT,
			Password:   ROOT_PASSWORD,
			RoleID:     role.ID,
			IsReserved: 1,
		}
		if CreateUser(&user) == nil {
			println("User root initialized")
		}

		_, err := InitializeUserSpaceFiles(user)
		if err != nil {
			println(err.Error())
		}
	}
}
