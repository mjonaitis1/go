package user

// ListUsers returns a slice of all available user entries.
func ListUsers() ([]*User, error) {
	return allUsers()
}

// ListGroups returns a slice of all available group entries.
func ListGroups() ([]*Group, error) {
	return allGroups()
}
