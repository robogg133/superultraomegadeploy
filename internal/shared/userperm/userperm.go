package userperm

const (
	CanChangeServerSettings string = "canChangeServerSettings"
	CanCreateOrg            string = "canCreateOrg"
	CanInviteUser           string = "canInviteUser"
)

// All lists every permission a user can hold.
var All = []string{CanCreateOrg, CanInviteUser, CanChangeServerSettings}
