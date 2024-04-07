package crud

// TeamCreate struct to describe create team objects.
type TeamCreate struct {
	TeamName string `json:"team_name" validate:"required,lte=255"`
	UserID   string `json:"user_id" validate:"required,lte=255"`
}

// TeamAddMember struct to describe add member to team objects.
type TeamAddMember struct {
	TeamID   string `json:"team_id" validate:"required,lte=255"`
	MemberID string `json:"member_id" validate:"required,lte=255"`
	UserID   string `json:"user_id" validate:"required,lte=255"`
}

// TeamRemoveMember struct to describe remove member from team objects.
type TeamRemoveMember struct {
	TeamID   string `json:"team_id" validate:"required,lte=255"`
	MemberID string `json:"member_id" validate:"required,lte=255"`
	UserID   string `json:"user_id" validate:"required,lte=255"`
}

// TeamListMembers struct to describe list members of team objects.
type TeamListMembers struct {
	TeamID string `json:"team_id" validate:"required,lte=255"`
	UserID string `json:"user_id" validate:"required,lte=255"`
}
