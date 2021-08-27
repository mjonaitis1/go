//go:build (aix || dragonfly || (!android && linux) || solaris) && cgo && !osusergo
// +build aix dragonfly !android,linux solaris
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

static FILE* openUsersFile(){
	FILE* fp;
	fp = fopen("./testdata/users.txt", "r");
	return fp;
}

static FILE* openGroupsFile(){
	FILE* fp;
	fp = fopen("./testdata/groups.txt", "r");
	return fp;
}
*/
import "C"

import (
	"os"
)

// iterateUsersTest implements usersIterator interface and is used for testing
// users iteration functionality with fgetpwent(3).
type iterateUsersTest struct {
	fp *C.FILE
}

func (i *iterateUsersTest) set() {
	var fp *C.FILE
	C.resetErrno()
	fp, err := C.openUsersFile()
	if err != nil {
		panic(err)
	}
	i.fp = fp
}

func (i *iterateUsersTest) get() (*C.struct_passwd, error) {
	var result *C.struct_passwd
	// fgetpwent(3) returns ENOENT when there are no more records. This is
	// undocumented in fgetgrent documentation, however, underlying
	// implementation of fgetpwent uses fgetpwent_r(3), which in turn returns
	// ENOENT when there are no more records.
	result, err := C.fgetpwent(i.fp)
	return result, err
}

func (i *iterateUsersTest) end() {
	if i.fp != nil {
		C.fclose(i.fp)
	}
}

// iterateGroupsTest implements groupsIterator interface and is used for testing
// users iteration functionality with fgetgrent(3).
type iterateGroupsTest struct {
	f  *os.File
	fp *C.FILE
}

func (i *iterateGroupsTest) set() {
	var fp *C.FILE
	C.resetErrno()
	fp, err := C.openGroupsFile()
	if err != nil {
		panic(err)
	}
	i.fp = fp
}

func (i *iterateGroupsTest) get() (*C.struct_group, error) {
	var result *C.struct_group
	result, err := C.fgetgrent(i.fp)
	// fgetgrent(3) returns ENOENT when there are no more records. This is
	// undocumented in fgetgrent documentation, however, underlying
	// implementation of fgetgrent uses fgetgrent_r(3), which in turn returns
	// ENOENT when there are no more records.
	return result, err
}

func (i *iterateGroupsTest) end() {
	if i.fp != nil {
		C.fclose(i.fp)
	}
}
