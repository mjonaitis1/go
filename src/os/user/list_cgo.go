//go:build (aix || dragonfly || freebsd || (!android && linux) || netbsd || openbsd || solaris || darwin) && cgo && !osusergo
// +build aix dragonfly freebsd !android,linux netbsd openbsd solaris darwin
// +build cgo
// +build !osusergo

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

// usersIterator defines the methods used in users iteration process within
// allUsers function. This interface allows testing allUsers functionality.
// list_test_fgetent.go file defines test related struct that implements
// usersIterator.
type usersIterator interface {
	// set sets up internal state before iteration
	set()

	// get sequentially returns a passwd structure which is later processed into
	// *User entry
	get() (*C.struct_passwd, error)

	// end cleans up internal state after iteration is done
	end()
}

type iterateUsers struct{}

func (i iterateUsers) set() {
	C.setpwent()
}

func (i iterateUsers) get() (*C.struct_passwd, error) {
	var result *C.struct_passwd
	result, err := C.getpwent()
	return result, err
}

func (i iterateUsers) end() {
	C.endpwent()
}

// This helper is used to retrieve users via c library call. A global
// variable which implements usersIterator interface is needed in order to
// separate testing logic from production. Since cgo can not be used directly
// in tests, list_test_fgetent.go file provides iterateUsersTest
// structure which implements usersIterator interface and can substitute
// default userIterator value.
var userIterator usersIterator = iterateUsers{}

// allUsers collects slice of users via getpwent(3). If error occurs during
// getpwent call, users slice is returned along with error.
//
// Since iterateUsers uses getpwent(3), which is not thread safe, allUsers
// can not bet used concurrently. If concurrent usage is required, it is
// recommended to use locking mechanism such as sync.Mutex when calling
// iterateUsers from multiple goroutines.
func allUsers() (users []*User, err error) {
	userIterator.set()
	defer userIterator.end()
	for {
		var result *C.struct_passwd
		C.resetErrno()
		result, err = userIterator.get()

		// If result is nil - getpwent iterated through entire users database
		// or there was an error
		if result == nil || err != nil {
			return
		}
		users = append(users, buildUser(result))
	}
}

// groupsIterator defines the methods used in groups iteration process within
// allGroups. This interface allows testing allGroups functionality.
// list_test_fgetent.go file defines test related struct that implements groupsIterator.
type groupsIterator interface {
	// set sets up internal state before iteration
	set()

	// get sequentially returns a group structure which is later processed into *Group entry
	get() (*C.struct_group, error)

	// end cleans up internal state after iteration is done
	end()
}

type iterateGroups struct{}

func (i iterateGroups) set() {
	C.setgrent()
}

func (i iterateGroups) get() (*C.struct_group, error) {
	var result *C.struct_group
	result, err := C.getgrent()
	return result, err
}

func (i iterateGroups) end() {
	C.endgrent()
}

// This helper is used to retrieve groups via c library call. A global
// variable which implements groupsIterator interface is needed in order to
// separate testing logic from production. Since cgo can not be used directly
// in tests, list_test_fgetent.go file provides iterateGroupsTest
// structure which implements groupsIterator interface and can substitute
// default groupIterator value.
var groupIterator groupsIterator = iterateGroups{}

// allGroups collects slice of groups via getgrent(3). If error occurs during
// getgrent call, groups slice is returned along with error.
//
// Since iterateGroups uses getgrent(3), which is not thread safe, iterateGroups
// can not bet used concurrently. If concurrent usage is required, it is
// recommended to use locking mechanism such as sync.Mutex when calling
// iterateGroups from multiple goroutines.
func allGroups() (groups []*Group, err error) {
	groupIterator.set()
	defer groupIterator.end()
	for {
		var result *C.struct_group
		C.resetErrno()
		result, err = groupIterator.get()

		// If result is nil - getgrent iterated through entire groups database
		// or there was an error
		if result == nil || err != nil {
			return
		}
		groups = append(groups, buildGroup(result))
	}
}
