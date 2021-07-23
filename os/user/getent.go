package user

func GetAllUsers() ([]*User, error){
	return allUsers()
}

func GetAllGroups() ([]*Group, error){
	return allGroups()
}