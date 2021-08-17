package models

func initDefaultRoles() {
	role := Role{
		Name:  ROOT_ROLENAME,
		Space: storage.Available,
		Code:  generateToken(ROLE_CODE_SIZE),
	}

	if db.Create(&role).Error == nil {
		println("Default roles initialized")
	}
}
