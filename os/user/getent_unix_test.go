// +build aix darwin dragonfly freebsd js,wasm !android,linux netbsd openbsd solaris
// +build !cgo osusergo

package user

import (
	"strings"
	"testing"
)

func TestIterateGroups(t *testing.T) {
	const testGroupFile = `# See the opendirectoryd(8) man page for additional 
# information about Open Directory.
##
nobody:*:-2:
nogroup:*:-1:
invalidgid:*:notanumber:root
+plussign:*:20:root
    indented:*:7:
# comment:*:4:found
     # comment:*:4:found
kmem:*:2:root
`
	// Ordered list of correctly parsed groups from testGroupFile
	var groups = []*Group{
		{
			Gid:  "-2",
			Name: "nobody",
		},
		{
			Gid:  "-1",
			Name: "nogroup",
		},
		{
			Gid:  "7",
			Name: "indented",
		},
		{
			Gid:  "2",
			Name: "kmem",
		},
	}

	r := strings.NewReader(testGroupFile)

	collectedGroups := make([]*Group, 0, 4)
	_, err := readColonFile(r, groupsIterator(func(g *Group) error {
		collectedGroups = append(collectedGroups, g)
		return nil
	}))

	if len(collectedGroups) != len(groups) {
		t.Errorf("groups were not parsed correctly: parsed %d/%d", len(collectedGroups), len(groups))
	}

	for i, g := range collectedGroups {
		if g.Name != groups[i].Name || g.Gid != groups[i].Gid {
			t.Errorf("groups could not be parsed correctly: parsed: %+v, actual: %+v", g, groups[i])
		}
	}

	if err != nil {
		t.Errorf("readEtcFile error: %v", err)
	}
}

func TestIterateUsers(t *testing.T) {
	const testUserFile = `   # Example user file
root:x:0:0:root:/root:/bin/bash
     indented:x:3:3:indented with a name:/dev:/usr/sbin/nologin
negative:x:-5:60:games:/usr/games:/usr/sbin/nologin
allfields:x:6:12:mansplit,man2,man3,man4:/home/allfields:/usr/sbin/nologin
+plussign:x:8:10:man:/var/cache/man:/usr/sbin/nologin

malformed:x:27:12 # more:colons:after:comment

struid:x:notanumber:12 # more:colons:after:comment

# commented:x:28:12:commented:/var/cache/man:/usr/sbin/nologin
      # commentindented:x:29:12:commentindented:/var/cache/man:/usr/sbin/nologin

struid2:x:30:badgid:struid2name:/home/struid:/usr/sbin/nologin
`
	var users = []*User{
		{
			Username: "root",
			Name:     "root",
			Uid:      "0",
			Gid:      "0",
			HomeDir:  "/root",
		},
		{
			Username: "indented",
			Name:     "indented with a name",
			Uid:      "3",
			Gid:      "3",
			HomeDir:  "/dev",
		},
		{
			Username: "negative",
			Name:     "games",
			Uid:      "-5",
			Gid:      "60",
			HomeDir:  "/usr/games",
		},
		{
			Username: "allfields",
			Name:     "mansplit",
			Uid:      "6",
			Gid:      "12",
			HomeDir:  "/home/allfields",
		},
	}

	collectedUsers := make([]*User, 0, 4)
	r := strings.NewReader(testUserFile)
	_, err := readColonFile(r, usersIterator(func(u *User) error {
		collectedUsers = append(collectedUsers, u)
		return nil
	}))

	if len(collectedUsers) != len(users) {
		t.Errorf("users were not parsed correctly: parsed %d/%d", len(collectedUsers), len(users))
	}

	for i, u := range collectedUsers {
		if u.Name != users[i].Name ||
			u.Gid != users[i].Gid ||
			u.Uid != users[i].Uid ||
			u.Username != users[i].Username ||
			u.HomeDir != users[i].HomeDir {
			t.Errorf("users could not be parsed correctly: parsed: %+v, actual: %+v", u, users[i])
		}
	}

	if err != nil {
		t.Errorf("readEtcFile error: %v", err)
	}
}
