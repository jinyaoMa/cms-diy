package model

func prepareDefaultRoles() {
	_, hasRole := GetRoleByName(ROOT_ROLENAME)
	if !hasRole {
		role := Role{
			Name:  ROOT_ROLENAME,
			Space: storage.Available,
		}

		if db.Create(&role).Error == nil {
			println("Default roles initialized")
		}
	}
}
