package user

import (
	"fmt"
	"internal/syscall/windows/registry"
	"syscall"
)

// ProfileList registry key defines all local users
const ProfileListKey = `SOFTWARE\Microsoft\Windows NT\CurrentVersion\ProfileList`

// ??
const HiveListKey = `Computer\HKEY_LOCAL_MACHINE\SYSTEM\CurrentControlSet\Control\hivelist`

// List of SIDs (sub keys of ProfileList)
func profileListSIDs() ([]string, error) {
	key, err := registry.OpenKey(registry.LOCAL_MACHINE, ProfileListKey, registry.QUERY_VALUE|registry.ENUMERATE_SUB_KEYS)
	if err != nil {
		return nil, err
	}
	return key.ReadSubKeyNames()
}

func iterateUsers(fn NextUserFunc) error {
	sids, err := profileListSIDs()
	if err != nil {
		return err
	}

	// If error occurs - we want to ignore it and continue with next entry in sids slice.
	for _, sidString := range sids {
		sid, e := syscall.StringToSid(sidString)
		if e != nil {
			continue
		}

		// Skip non user accounts
		if _, _, accType, _ := sid.LookupAccount(""); accType != syscall.SidTypeUser {
			continue
		}

		u, err := newUserFromSid(sid)
		if err != nil {
			continue
		}

		// Callback to user supplied fn, with user
		if err := fn(u); err != nil {
			return err
		}
	}
	return nil
}

func iterateGroups(fn NextGroupFunc) error {
	fmt.Println("windows iterateUser")
	return nil
}
