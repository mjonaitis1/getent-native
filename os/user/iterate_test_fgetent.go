// +build aix dragonfly !android,linux solaris
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

import (
	"errors"
	"os"
	"syscall"
	"unsafe"
)

const (
	_usersTestDataFile  = "./testdata/users.txt"
	_groupsTestDataFile = "./testdata/groups.txt"
)

// iterateUsersHelperTest implements usersHelper interface and is used for testing
// users iteration functionality with fgetpwent library call.
type iterateUsersHelperTest struct {
	f  *os.File
	fp *C.FILE
}

func (i *iterateUsersHelperTest) set() {
	mode := C.CString("r")
	defer C.free(unsafe.Pointer(mode))
	f, err := os.Open(_usersTestDataFile)
	if err != nil {
		panic(err)
	}
	i.f = f
	var fp *C.FILE
	C.resetErrno()
	fp, err = C.fdopen(C.int(f.Fd()), mode)
	if fp == nil {
		panic(err)
	}
	i.fp = fp
}

func (i *iterateUsersHelperTest) get() (*C.struct_passwd, error) {
	var result *C.struct_passwd
	result, err := C.fgetpwent(i.fp)

	// fgetpwent returns ENOENT when there are no more records
	if result == nil && errors.Is(err, syscall.ENOENT) {
		return nil, nil
	}

	return result, err
}

func (i *iterateUsersHelperTest) end() {
	if i.f != nil {
		_ = i.f.Close()
	}
}

// iterateGroupsHelperTest implements groupsHelper interface and is used for testing
// users iteration functionality with fgetgrent library call.
type iterateGroupsHelperTest struct {
	f  *os.File
	fp *C.FILE
}

func (i *iterateGroupsHelperTest) set() {
	mode := C.CString("r")
	defer C.free(unsafe.Pointer(mode))
	f, err := os.Open(_groupsTestDataFile)
	if err != nil {
		panic(err)
	}
	i.f = f
	var fp *C.FILE
	C.resetErrno()
	fp, err = C.fdopen(C.int(f.Fd()), mode)
	if fp == nil {
		panic(err)
	}
	i.fp = fp
}

func (i *iterateGroupsHelperTest) get() (*C.struct_group, error) {
	var result *C.struct_group
	result, err := C.fgetgrent(i.fp)

	// fgetpwent returns ENOENT when there are no more records
	if result == nil && errors.Is(err, syscall.ENOENT) {
		return nil, nil
	}

	return result, err
}

func (i *iterateGroupsHelperTest) end() {
	if i.f != nil {
		_ = i.f.Close()
	}
}
