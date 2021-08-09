// +build aix dragonfly freebsd !android,linux netbsd openbsd solaris darwin
// +build cgo,!osusergo

package user

/*
#include <unistd.h>
#include <sys/types.h>
#include <pwd.h>
#include <grp.h>
#include <stdlib.h>
#include <stdio.h>
#include <errno.h>

static void resetErrno(){
	errno = 0;
}
*/
import "C"

// usersHelper defines the methods used in users iteration process within iterateUsers. This interface allows to test
// iterateUsers functionality. iterate_test_fgetent.go file defines test related struct that implements usersHelper.
type usersHelper interface {
	// set does any necessary setup steps before starting iteration
	set()

	// get sequentially returns a passwd structure which is later processed into *User entry
	get() (*C.struct_passwd, error)

	// end does any necessary finalization steps after iteration is done
	end()
}

type iterateUsersHelper struct{}

func (i iterateUsersHelper) set() {
	C.setpwent()
}

func (i iterateUsersHelper) get() (*C.struct_passwd, error) {
	var result *C.struct_passwd
	result, err := C.getpwent()
	return result, err
}

func (i iterateUsersHelper) end() {
	C.endpwent()
}

// iterate users helper
var iuh usersHelper = iterateUsersHelper{}

// iterateUsers iterates over users database via getgrent library call. If fn returns non nil error, then
// iteration is terminated. A nil result from getgrent means there were no more entries, or an error occurred,
// as such, iteration is terminated, and if error was encountered it is returned.
//
// Since iterateUsers uses getgrent library call, which is not thread safe, iterateUsers can not bet used concurrently.
// If concurrent usage is required, it is recommended to use locking mechanism such as sync.Mutex when calling
// iterateUsers from multiple goroutines.
func iterateUsers(fn NextUserFunc) error {
	iuh.set()
	defer iuh.end()
	for {
		var result *C.struct_passwd
		C.resetErrno()
		result, err := iuh.get()

		// If result is nil - getgrent iterated through entire users database or there was an error
		if result == nil {
			return err
		}

		// User provided non-nil error means that iteration should be terminated
		if err = fn(buildUser(result)); err != nil {
			return err
		}
	}
}

// groupsHelper defines the methods used in groups iteration process within iterateGroups. This interface allows to test
// iterateGroups functionality. iterate_test_fgetent.go file defines test related struct that implements groupsHelper.
type groupsHelper interface {
	// set does any necessary setup steps before starting iteration
	set()

	// get sequentially returns a group structure which is later processed into *Group entry
	get() (*C.struct_group, error)

	// end does any necessary finalization steps after iteration is done
	end()
}

type iterateGroupsHelper struct{}

func (i iterateGroupsHelper) set() {
	C.setgrent()
}

func (i iterateGroupsHelper) get() (*C.struct_group, error) {
	var result *C.struct_group
	result, err := C.getgrent()
	return result, err
}

func (i iterateGroupsHelper) end() {
	C.endgrent()
}

// iterate users helper
var igh groupsHelper = iterateGroupsHelper{}

// iterateGroups iterates over groups database via getgrent library call. If fn returns non nil error, then
// iteration is terminated. A nil result from getgrent means there were no more entries, or an error occurred,
// as such, iteration is terminated, and if error was encountered it is returned.
//
// Since iterateGroups uses getgrent library call, which is not thread safe, iterateGroups can not bet used concurrently.
// If concurrent usage is required, it is recommended to use locking mechanism such as sync.Mutex when calling
// iterateGroups from multiple goroutines.
func iterateGroups(fn NextGroupFunc) error {
	igh.set()
	defer igh.end()
	for {
		var result *C.struct_group
		C.resetErrno()
		result, err := igh.get()

		// If result is nil - getgrent iterated through entire groups database or there was an error
		if result == nil {
			return err
		}

		// User provided non-nil error means that iteration should be terminated
		if err = fn(buildGroup(result)); err != nil {
			return err
		}
	}
}
