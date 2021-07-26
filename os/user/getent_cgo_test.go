// +build aix dragonfly freebsd !android,linux netbsd openbsd solaris darwin
// +build cgo,!osusergo

package user

import (
	"errors"
	"testing"
)

// As we can't really test cgo, we will attempt to
// check if user/group record could be at least retrieved

func TestIterateUser(t *testing.T) {
	err := iterateUsers(func(user *User) error {
		if user.Username == "" && user.Gid == "" && user.Uid == "" && user.HomeDir == "" && user.Name == "" {
			t.Errorf("parsed user is empty: %+v \n", user)
		}
		return errors.New("terminate iteration")
	})

	if err != nil {
		t.Errorf("error encoutered while iterating users: %w", err)
	}
}

func TestIterateGroup(t *testing.T) {
	err := iterateGroups(func(group *Group) error {
		if group.Name == "" && group.Gid == "" {
			t.Errorf("parsed group is empty: %+v \n", group)
		}
		return errors.New("terminate iteration")
	})

	if err != nil {
		t.Errorf("error encoutered while iterating groups: %w", err)
	}
}
