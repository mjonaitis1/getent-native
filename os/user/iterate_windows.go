package user

import (
	"internal/syscall/windows/registry"
	"syscall"
)

// ProfileListKey registry key contains all local users SIDs as sub keys. It is a sub key of HKEY_LOCAL_MACHINE.
const ProfileListKey = `SOFTWARE\Microsoft\Windows NT\CurrentVersion\ProfileList`

// iterateSIDS iterates through ProfileListKey sub keys and calls provided fn with enumerated sub key
// name as parameter. If fn returns non-nil error, iteration is terminated.
// This function is modified version of registry.Key.ReadSubKeyNames.
func iterateSIDS(fn func(string) error) error {
	k, err := registry.OpenKey(registry.LOCAL_MACHINE, ProfileListKey, registry.QUERY_VALUE|registry.ENUMERATE_SUB_KEYS)
	if err != nil {
		return err
	}

	// Registry key size limit is 255 bytes and described there:
	// https://msdn.microsoft.com/library/windows/desktop/ms724872.aspx
	buf := make([]uint16, 256) //plus extra room for terminating zero byte
loopItems:
	for i := uint32(0); ; i++ {
		l := uint32(len(buf))
		for {
			err := syscall.RegEnumKeyEx(syscall.Handle(k), i, &buf[0], &l, nil, nil, nil, nil)
			if err == nil {
				break
			}
			if err == syscall.ERROR_MORE_DATA {
				// Double buffer size and try again.
				l = uint32(2 * len(buf))
				buf = make([]uint16, l)
				continue
			}
			// ERROR_NO_MORE_ITEMS means all keys have been iterated.
			if err == syscall.Errno(259) {
				break loopItems
			}
		}

		// Call fn and provide the sub key name, which is a SID string.
		sid := syscall.UTF16ToString(buf[:l])
		if err := fn(sid); err != nil {
			return err
		}
	}

	return nil
}

// iterateUsers iterates through ProfileListKey SIDs, looks up for user with each given SID and calls user provided fn
// with each *User entry. Each iterated SID can be either user or group. Only user SIDs are processed.
func iterateUsers(fn NextUserFunc) error {
	return iterateSIDS(func(sid string) error {
		SID, err := syscall.StringToSid(sid)

		if err != nil {
			return nil
		}
		// Skip non user SID
		if _, _, accType, _ := SID.LookupAccount(""); accType != syscall.SidTypeUser {
			return nil
		}
		u, err := newUserFromSid(SID)
		if err != nil {
			return nil
		}

		// Callback to user supplied fn, with user
		if err := fn(u); err != nil {
			return err
		}

		return nil
	})
}

// iterateGroups iterates through ProfileListKey SIDs, looks up for group with each given SID and calls user provided fn
// with each *Group entry. Each iterated SID can be either user or group. Only group SIDs are processed.
func iterateGroups(fn NextGroupFunc) error {
	return iterateSIDS(func(sid string) error {
		SID, err := syscall.StringToSid(sid)
		if err != nil {
			return nil
		}

		groupname, _, t, err := SID.LookupAccount("")
		if err != nil {
			return err
		}
		// Skip non groups
		if t != syscall.SidTypeGroup && t != syscall.SidTypeWellKnownGroup && t != syscall.SidTypeAlias {
			return nil
		}
		g := &Group{Name: groupname, Gid: sid}

		// Callback to user supplied fn, with group
		if err := fn(g); err != nil {
			return err
		}

		return nil
	})
}
