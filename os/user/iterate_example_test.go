package user

import (
	"errors"
	"fmt"
)

func ExampleIterateUsers() {
	// Get first 20 users
	users := make([]*User, 0, 20)
	i := 0
	err := IterateUsers(func(user *User) error {
		users = append(users, user)
		i++

		// Once we return non-nil error - iteration process stops
		if i >= 20 {
			return errors.New("stop iterating")
		}

		// As long as error is nil, IterateUsers will iterate over users database
		return nil
	})

	if err != nil {
		fmt.Printf("error encountered while iterating users database: %v", err)
	}

	// Do something with users slice
}

func ExampleIterateGroups() {
	// Get first 20 groups
	groups := make([]*Group, 0, 20)
	i := 0
	err := IterateGroups(func(group *Group) error {
		groups = append(groups, group)
		i++

		// Once we return non-nil error - iteration process stops
		if i >= 20 {
			return errors.New("stop iterating")
		}

		// As long as error is nil, IterateGroups will iterate over groups database
		return nil
	})

	if err != nil {
		fmt.Printf("error encountered while iterating groups database: %v", err)
	}

	// Do something with groups slice
}
