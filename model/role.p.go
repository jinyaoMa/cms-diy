package model

func prepareDefaultRoles() {
	_, hasRole := GetRoleByName(ROOT_ROLENAME)
	if !hasRole {
		roles := Roles{
			{
				Name:       ROOT_ROLENAME,
				Space:      storage.Available,
				Code:       generateToken(ROLE_CODE_SIZE),
				IsReserved: 1,
			},
			{
				Name:       ROLE_DEFAULT_MEMBER_NAME,
				Permission: ROLE_DEFAULT_MEMBER_PERMISSION,
				Space:      ROLE_DEFAULT_MEMBER_SPACE,
				Code:       ROLE_DEFAULT_MEMBER_CODE,
				IsReserved: 1,
			},
		}

		if db.Create(&roles).Error == nil {
			println("Default roles initialized")
		}
	}
}
