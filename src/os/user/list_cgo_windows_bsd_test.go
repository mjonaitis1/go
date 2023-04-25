//go:build ((darwin || freebsd || openbsd || netbsd) && cgo && !osusergo) || windows

package user

import (
	"testing"
)

// As BSDs (including darwin) do not support fgetpwent(3)/fgetgrent(3), attempt
// to check if user/group records at least can be retrieved.
// On Windows, it is not possible to easily mock registry. Checking if
// user/group records can at least be retrieved will suffice.

func TestIterateUser(t *testing.T) {
	users, err := allUsers()
	if err != nil {
		t.Errorf("iterating users: %w", err)
	}
	if len(users) == 0 {
		t.Errorf("no users were retrieved")
	}
}

func TestIterateGroup(t *testing.T) {
	groups, err := allGroups()
	if err != nil {
		t.Errorf("iterating groups: %w", err)
	}
	if len(groups) == 0 {
		t.Errorf("no groups were retrieved")
	}
}
