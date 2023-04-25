package user

import (
	"internal/syscall/windows/registry"
	"syscall"
)

// _profileListKey registry key contains all local user/group SIDs
// (Security Identifiers are Windows version of user/group ids on unix systems)
// as sub keys. It is a sub key of HKEY_LOCAL_MACHINE. Since
// HKEY_LOCAL_MACHINE\SOFTWARE\Microsoft\Windows NT\CurrentVersion\ProfileList
// registry key does not contain SIDs of Administrator, Guest, or other default
// system users, some users or groups might not be provided when using
// allUsers or allGroups
const _profileListKey = `SOFTWARE\Microsoft\Windows NT\CurrentVersion\ProfileList`

// allSIDS returns a slice of _profileListKey sub key names, which are
// essentially SIDs.
func allSIDS() ([]string, error) {
	k, err := registry.OpenKey(registry.LOCAL_MACHINE, _profileListKey, registry.QUERY_VALUE|registry.ENUMERATE_SUB_KEYS)
	if err != nil {
		return nil, err
	}
	return k.ReadSubKeyNames()
}

// allUsers loops through SIDs, looks up for user with each given SID
// and if a user was found, appends it to users slice.
func allUsers() ([]*User, error) {
	users := make([]*User, 0)
	sids, err := allSIDS()
	if err != nil {
		return nil, err
	}
	for _, sid := range sids {
		SID, err := syscall.StringToSid(sid)
		if err != nil {
			return users, err
		}

		// Skip non-user SID
		if _, _, accType, _ := SID.LookupAccount(""); accType != syscall.SidTypeUser {
			continue
		}
		u, err := newUserFromSid(SID)
		if err != nil {
			return users, err
		}
		users = append(users, u)
	}
	return users, err
}

// allGroups loops through SIDs, looks up for group with each given SID
// and if a group was found, appends it to groups slice.
func allGroups() ([]*Group, error) {
	groups := make([]*Group, 0)
	sids, err := allSIDS()
	if err != nil {
		return nil, err
	}
	for _, sid := range sids {
		SID, err := syscall.StringToSid(sid)
		if err != nil {
			return groups, err
		}

		groupname, _, t, err := SID.LookupAccount("")
		if err != nil {
			return groups, err
		}
		// Skip non-group SID
		if isNotGroup(t) {
			continue
		}
		g := &Group{Name: groupname, Gid: sid}

		// Callback to user supplied fn, with group
		groups = append(groups, g)
	}
	return groups, err
}
