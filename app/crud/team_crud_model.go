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
