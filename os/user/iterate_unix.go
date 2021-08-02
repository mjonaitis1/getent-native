// +build aix darwin dragonfly freebsd js,wasm !android,linux netbsd openbsd solaris
// +build !cgo osusergo

package user

import (
	"os"
)

func iterateUsers(fn NextUserFunc) error {
	f, err := os.Open(userFile)
	if err != nil {
		return err
	}
	defer f.Close()
	_, err = readColonFile(f, usersIterator(fn))
	return err
}

func iterateGroups(fn NextGroupFunc) error {
	f, err := os.Open(groupFile)
	if err != nil {
		return err
	}
	defer f.Close()
	_, err = readColonFile(f, groupsIterator(fn))
	return err
}

// usersIterator does not return value other than error, as we want to make sure that we iterate through
// entire file with readColonFile. If error is returned from fn - iteration process is terminated
func usersIterator(fn NextUserFunc) lineFunc {
	return func(line []byte) (interface{}, error) {
		v, _ := matchUserIndexValue("", -1)(line)
		if u, ok := v.(*User); ok {
			err := fn(u)
			if err != nil {
				return nil, err
			}
		}
		return nil, nil
	}
}

// groupsIterator does not return value other than error, as we want to make sure that we iterate through
// entire file with readColonFile. If error is returned from fn - iteration process is terminated
func groupsIterator(fn NextGroupFunc) lineFunc {
	return func(line []byte) (interface{}, error) {
		v, _ := matchGroupIndexValue("", -1)(line)
		if g, ok := v.(*Group); ok {
			err := fn(g)
			if err != nil {
				return nil, err
			}
		}
		return nil, nil
	}
}
