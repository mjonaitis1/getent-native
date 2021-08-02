// +build darwin freebsd openbsd netbsd
// +build cgo,!osusergo

package user

import (
	"errors"
	"testing"
)

// As we can't really test cgo implementation on darwin or bsd based oses, since there is no support for fgetpwent/fgetgrent
// library calls, we will attempt to check if user/group record could be at least retrieved

var _stopErr = errors.New("terminate iteration")

func TestIterateUser(t *testing.T) {
	err := iterateUsers(func(user *User) error {
		if user.Username == "" && user.Gid == "" && user.Uid == "" && user.HomeDir == "" && user.Name == "" {
			t.Errorf("parsed user is empty: %+v \n", user)
		}
		return _stopErr
	})

	if !errors.Is(err, _stopErr) {
		t.Errorf("error encoutered while iterating users: %w", err)
	}
}

func TestIterateGroup(t *testing.T) {
	err := iterateGroups(func(group *Group) error {
		if group.Name == "" && group.Gid == "" {
			t.Errorf("parsed group is empty: %+v \n", group)
		}
		return _stopErr
	})

	if !errors.Is(err, _stopErr) {
		t.Errorf("error encoutered while iterating groups: %w", err)
	}
}
