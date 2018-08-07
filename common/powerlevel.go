package common

import (
	"errors"
	"fmt"
)

const (
	// Room visibility
	RoomVisibilityPrivate = "private"
	RoomVisibilityPublic = "public"

	// User/team power levels
	PowerlevelUser = 0
	PowerlevelModerator = 50
	PowerlevelAdmin = 100

	// Actions
	ActionDeleteRoom = "DeleteRoom"
	ActionSetStateRoom = "SetStateRoom"
	ActionRedactionRoom = "RedactionRoom"
	ActionKickRoom = "KickRoom"
	ActionBanRoom = "BanRoom"

	ActionAddUserOrg = "AddUserOrg"
	ActionCreateTeamOrg = "CreateTeamOrg"
	ActionUpdateTeamOrg = "UpdateTeamOrg"
	ActionDeleteTeamOrg = "DeleteTeamOrg"
	ActionAddUserTeamOrg = "AddUserTeamOrg"
	ActionListTeamsOrg = "ListTeamsOrg"
	ActionListMembersTeamsOrg = "ListMembersTeamsOrg"
	ActionDeleteMembersTeamsOrg = "DeleteMembersTeamsOrg"
)

var (
	MinimumPowerLevel map[string]int64 = map[string]int64{
		ActionDeleteRoom: PowerlevelAdmin,
		ActionSetStateRoom: PowerlevelAdmin,
		ActionRedactionRoom: PowerlevelAdmin,
		ActionKickRoom: PowerlevelAdmin,
		ActionBanRoom: PowerlevelAdmin,

		ActionAddUserOrg : PowerlevelAdmin,
		ActionCreateTeamOrg : PowerlevelAdmin,
		ActionUpdateTeamOrg : PowerlevelAdmin,
		ActionDeleteTeamOrg : PowerlevelAdmin,
		ActionAddUserTeamOrg : PowerlevelAdmin,
		ActionListTeamsOrg : PowerlevelAdmin,
		ActionListMembersTeamsOrg : PowerlevelAdmin,
		ActionDeleteMembersTeamsOrg : PowerlevelAdmin,
	}
)


// Checks if this power level is enough to perform given action
func EnoughPowerLevel(action string, level int64) (bool, error) {
	levelNeeded, ok := MinimumPowerLevel[action]
	if !ok {
		return false, errors.New(fmt.Sprintf("Action not found %v", action))
	}
	if levelNeeded > level {
		return false, errors.New(fmt.Sprintf("Not enough power level to perform %v", action))
	}
	return true, nil
}