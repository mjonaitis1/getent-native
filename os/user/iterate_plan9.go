package user

import (
	"bytes"
	"fmt"
	"os"
	"strings"
)

// users and groups file location in plan9
var (
	usersFile = "/adm/users"
)

// userGroupIterator is a helper iterator function, which parses /adm/users
func userGroupIterator(lineFn lineFunc) error {
	f, err := os.Open(usersFile)
	if err != nil {
		return fmt.Errorf("opening users file: %w", err)
	}
	defer func() {
		_ = f.Close()
	}()

	_, err = readColonFile(f, lineFn)
	return err
}

// matchPlan9UserGroup matches valid /adm/users line.
// Since plan9 /adm/users lines are both users and groups,
// we just need to make sure that the line structure is correct.
//
// idx should be set to -1 for iteration.
//
// id:name:leader:members
// sys:sys::glenda,mj <-- user/group without a leader, with 2 members glenda and mj
// mj:mj:: <-- user/group without a leader, without members
//
//
// According to plan 9 users(6): ids are arbitrary text strings, typically the same as name.
// In older Plan 9 file servers, ids are small decimal numbers.
//
func matchPlan9UserGroup(value string, idx int, returnUser bool) lineFunc {
	var leadColon string
	if idx > 0 {
		leadColon = ":"
	}
	substr := []byte(leadColon + value + ":")
	return func(line []byte) (v interface{}, err error) {
		if (idx != -1 && !bytes.Contains(line, substr)) || bytes.Count(line, []byte{':'}) < 3 {
			return
		}
		// id:name:leader:members
		parts := strings.SplitN(string(line), ":", 4)
		if len(parts) < 4 || parts[0] == "" || (idx != -1 && parts[idx] != value) ||
			parts[0][0] == '+' || parts[0][0] == '-' {
			return
		}

		// Since plan9 stores a user and a group as the same record, we return the requested type.
		if returnUser {
			return &User{
				Uid:      parts[0],
				Gid:      parts[0],
				Username: parts[1],
				Name:     parts[1],
				HomeDir:  "usr/" + parts[1],
			}, nil
		}
		return &Group{Name: parts[1], Gid: parts[0]}, nil
	}
}

func iterateUsers(fn NextUserFunc) error {
	return userGroupIterator(func(line []byte) (interface{}, error) {
		v, _ := matchPlan9UserGroup("", -1, true)(line)
		if user, ok := v.(*User); ok {
			err := fn(user)
			if err != nil {
				return nil, err
			}
		}
		return nil, nil
	})
}

func iterateGroups(fn NextGroupFunc) error {
	return userGroupIterator(func(line []byte) (interface{}, error) {
		v, _ := matchPlan9UserGroup("", -1, true)(line)
		if group, ok := v.(*Group); ok {
			err := fn(group)
			if err != nil {
				return nil, err
			}
		}
		return nil, nil
	})
}
