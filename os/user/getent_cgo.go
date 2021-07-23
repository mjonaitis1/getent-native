// +build aix darwin dragonfly freebsd !android,linux netbsd openbsd solaris
// +build cgo,!osusergo

package user
/**
// Use getpwent
 */
import "C"

func allUsers() ([]*User, error){

}

func allGroups() ([]*Group, error){

}