/*
Package user is an extension of os/user package. Main goal is to provide access to users and groups iteration functionality.

Users and groups iteration functionality might be useful in cases where a full or partial list of all available
users/groups is required in application. As original os/users lookup functionality, iteration functionality relies on
two internal implementations: one that does not use cgo and parses contents of /etc/passwd and /etc/group files
for users and groups respectively, and another one that does use cgo and utilizes POSIX compliant libc library
routines getpwent and getgrent for users and groups respectively.

Example usage:

// Get all user names in system
usernames := make([]string, 0)
user.IterateUsers(func(u *user.User) error {
	usernames = append(usernames, u.Username)
	return nil
})

// Get all group names in system
groupnames := make([]string, 0)
user.IterateGroups(func(g *user.Group) error {
	groupnames = append(groupnames, g.Name)
	return nil
})

*/
package user
