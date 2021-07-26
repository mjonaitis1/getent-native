// +build aix dragonfly freebsd !android,linux netbsd openbsd solaris darwin
// +build cgo,!osusergo

package user

/*
#include <unistd.h>
#include <sys/types.h>
#include <pwd.h>
#include <grp.h>
#include <stdlib.h>
#include <errno.h>

static void resetErrno(){
	errno = 0;
}
*/
import "C"

// iterateUsers iterates over users database via getpwent library call. If fn returns non nil error, then
// iteration is terminated. A nil result from getpwent means there were no more entries, or an error occurred,
// as such, iteration is terminated, and if error was encountered it is returned.
//
// Since iterateUsers uses getpwent library call, which is not thread safe, iterateUsers can not bet used concurrently.
// If concurrent usage is required, it is recommended to use locking mechanism such as sync.Mutex when calling
// iterateUsers from multiple goroutines.
func iterateUsers(fn NextUserFunc) error {
	C.setpwent()
	for {
		var result *C.struct_passwd
		C.resetErrno()
		result, errno := C.getpwent()

		// If result is nil - getpwent iterated through entire users database or there was an error
		if result == nil {
			return errno
		}

		// User provided non-nil error means that iteration should be terminated
		if err := fn(buildUser(result)); err != nil {
			break
		}
	}
	C.endpwent()
	return nil
}

// iterateGroups iterates over groups database via getgrent library call. If fn returns non nil error, then
// iteration is terminated. A nil result from getgrent means there were no more entries, or an error occurred,
// as such, iteration is terminated, and if error was encountered it is returned.
//
// Since iterateGroups uses getgrent library call, which is not thread safe, iterateGroups can not bet used concurrently.
// If concurrent usage is required, it is recommended to use locking mechanism such as sync.Mutex when calling
// iterateGroups from multiple goroutines.
func iterateGroups(fn NextGroupFunc) error {
	C.setgrent()
	for {
		var result *C.struct_group
		C.resetErrno()
		result, errno := C.getgrent()

		// If result is nil - getgrent iterated through entire groups database or there was an error
		if result == nil {
			return errno
		}

		// User provided non-nil error means that iteration should be terminated
		if err := fn(buildGroup(result)); err != nil {
			break
		}
	}
	C.endgrent()
	return nil
}
