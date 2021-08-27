package user_test

import (
	"fmt"
	"os/user"
)

func ExampleListUsers() {
	users, err := user.ListUsers()
	for _, u := range users {
		fmt.Printf("Username: %s\n", u.Username)
	}
	if err != nil {
		fmt.Printf("error encountered while iterating users database: %v", err)
	}
}

func ExampleListGroups() {
	groups, err := user.ListGroups()
	for _, g := range groups {
		fmt.Printf("Groupname: %s\n", g.Name)
	}
	if err != nil {
		fmt.Printf("error encountered while iterating groups database: %v", err)
	}
}
