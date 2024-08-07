package controllers

import (
	"github.com/UIdealist/UIdealist-Members-Ms/app/crud"
	"github.com/UIdealist/UIdealist-Members-Ms/app/models"
	"github.com/UIdealist/UIdealist-Members-Ms/pkg/repository"
	"github.com/UIdealist/UIdealist-Members-Ms/platform/database"

	accessconnectormodels "github.com/UIdealist/Uidealist-Access-Ms/app/crud"
	accessconnector "github.com/UIdealist/Uidealist-Access-Ms/connector"

	"github.com/gofiber/fiber/v2"
)

// CreateTeam method to create a new team
// @Description Create a new team
// @Summary Create a new team
// @Tags Team
// @Accept json
// @Produce json
// @Param data body crud.TeamCreate true "Team Data"
// @Success 201 {string} status "ok"
// @Router /v1/team [post]
func CreateTeam(c *fiber.Ctx) error {
	addTeamRequest := &crud.TeamCreate{}

	// Checking received data from JSON body.
	if err := c.BodyParser(addTeamRequest); err != nil {
		// Return status 400 and error message.
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": true,
			"code":  repository.INVALID_DATA,
			"msg":   err.Error(),
		})
	}

	// Create database connection.
	db := database.DB

	// Create a new team.
	team := &models.Team{
		Name: addTeamRequest.TeamName,
	}

	// Create a new team in database.
	err := db.Create(&team).Error
	if err != nil {
		// Return status 500 and database connection error.
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": true,
			"code":  repository.DATABASE_ERROR,
			"msg":   err.Error(),
		})
	}

	// Get user from database
	user := &models.User{}

	// Find user by ID
	err = db.Where(&models.User{ID: addTeamRequest.UserID}).First(&user).Error
	if err != nil {
		// Return status 500 and database connection error.
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": true,
			"code":  repository.USER_NOT_FOUND,
			"msg":   err.Error(),
		})
	}

	// Add this user as a member of the team
	teamMember := &models.TeamHasMember{
		TeamID:            team.ID,
		MemberID:          user.MemberID,
		TeamMemberIsOwner: true,
	}

	// Create a new team member in database.
	err = db.Create(&teamMember).Error
	if err != nil {
		// Return status 500 and database connection error.
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": true,
			"code":  repository.DATABASE_ERROR,
			"msg":   err.Error(),
		})
	}

	// Grant creator admin access to the team
	response, err := accessconnector.GrantAccess(
		&accessconnectormodels.AccessList[string]{
			Policies: []accessconnectormodels.Access[string]{
				accessconnectormodels.Access[string]{
					Subject: "member-" + user.MemberID,
					Object:  "team-" + team.ID,
					Action:  "admin",
				},
			},
		},
	)

	if err != nil || response.Error {
		// Return status 500 and database connection error.
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": true,
			"code":  repository.DATABASE_ERROR,
			"msg":   err.Error(),
		})
	}

	// Return status 201 and created team.
	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"error": false,
		"code":  repository.TEAM_CREATED,
		"msg":   nil,
		"team":  team,
	})
}

// AddTeamMember method to add a member to a team
// @Description Add a member to a team
// @Summary Add a member to a team
// @Tags Team
// @Accept json
// @Produce json
// @Param data body crud.TeamAddMember true "New member data"
// @Success 201 {string} status "ok"
// @Router /v1/team/members [post]
func AddTeamMember(c *fiber.Ctx) error {
	addTeamMember := &crud.TeamAddMember{}

	// Checking received data from JSON body.
	if err := c.BodyParser(addTeamMember); err != nil {
		// Return status 400 and error message.
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": true,
			"code":  repository.INVALID_DATA,
			"msg":   err.Error(),
		})
	}

	// Create database connection.
	db := database.DB

	// Get user from database
	user := &models.User{}

	// Find user by ID
	err := db.Where(&models.User{ID: addTeamMember.UserID}).First(&user).Error
	if err != nil {
		// Return status 500 and database connection error.
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": true,
			"code":  repository.USER_NOT_FOUND,
			"msg":   err.Error(),
		})
	}

	// Check if user is owner of the team
	teamMember := &models.TeamHasMember{}

	err = db.Where(&models.TeamHasMember{TeamID: addTeamMember.TeamID, MemberID: user.MemberID}).First(&teamMember).Error
	if err != nil {
		// Return status 500 and database connection error.
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": true,
			"code":  repository.DATABASE_ERROR,
			"msg":   err.Error(),
		})
	}

	if teamMember == nil {
		// Return status 500 and database connection error.
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": true,
			"code":  repository.TEAM_MEMBER_NOT_FOUND,
			"msg":   "User is not member of the team",
		})
	}

	// TODO: Query permissions microservice to check if user is owner of the team
	if !teamMember.TeamMemberIsOwner {
		// Return status 500 and database connection error.
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": true,
			"code":  repository.INVALID_DATA,
			"msg":   "User is not owner of the team",
		})
	}

	// Add this user as a member of the team
	teamMember = &models.TeamHasMember{
		TeamID:            addTeamMember.TeamID,
		MemberID:          addTeamMember.MemberID,
		TeamMemberIsOwner: false,
	}

	// Create a new team member in database.
	err = db.Create(&teamMember).Error
	if err != nil {
		// Return status 500 and database connection error.
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": true,
			"code":  repository.DATABASE_ERROR,
			"msg":   err.Error(),
		})
	}

	// Return status 201 and created team.
	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"error": false,
		"code":  repository.TEAM_MEMBER_ADDED,
		"msg":   nil,
	})
}

// RemoveTeamMember method to remove a member from a team
// @Description Remove a member from a team
// @Summary Create a new team
// @Tags Team
// @Accept json
// @Produce json
// @Param data body crud.TeamRemoveMember true "Old member data"
// @Success 200 {string} status "ok"
// @Router /v1/team/members [delete]
func RemoveTeamMember(c *fiber.Ctx) error {
	removeTeamMember := &crud.TeamRemoveMember{}

	// Checking received data from JSON body.
	if err := c.BodyParser(removeTeamMember); err != nil {
		// Return status 400 and error message.
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": true,
			"code":  repository.INVALID_DATA,
			"msg":   err.Error(),
		})
	}

	// Create database connection.
	db := database.DB

	// Get user from database
	user := &models.User{}

	// Find user by ID
	err := db.Where(&models.User{ID: removeTeamMember.UserID}).First(&user).Error
	if err != nil {
		// Return status 500 and database connection error.
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": true,
			"code":  repository.USER_NOT_FOUND,
			"msg":   err.Error(),
		})
	}

	// Check if user is owner of the team
	teamMember := &models.TeamHasMember{}

	err = db.Where(&models.TeamHasMember{TeamID: removeTeamMember.TeamID, MemberID: user.MemberID}).First(&teamMember).Error
	if err != nil {
		// Return status 500 and database connection error.
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": true,
			"code":  repository.DATABASE_ERROR,
			"msg":   err.Error(),
		})
	}

	if teamMember == nil {
		// Return status 500 and database connection error.
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": true,
			"code":  repository.TEAM_MEMBER_NOT_FOUND,
			"msg":   "User is not member of the team",
		})
	}

	// TODO: Query permissions microservice to check if user is owner of the team
	if !teamMember.TeamMemberIsOwner {
		// Return status 500 and database connection error.
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": true,
			"code":  repository.INVALID_DATA,
			"msg":   "User is not owner of the team",
		})
	}

	// Check if member is member of the team
	err = db.Where(&models.TeamHasMember{TeamID: removeTeamMember.TeamID, MemberID: removeTeamMember.MemberID}).First(&teamMember).Error
	if err != nil {
		// Return status 500 and database connection error.
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": true,
			"code":  repository.DATABASE_ERROR,
			"msg":   err.Error(),
		})
	}

	if teamMember == nil {
		// Return status 500 and database connection error.
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": true,
			"code":  repository.INVALID_DATA,
			"msg":   "Specified member is not member of the team",
		})
	}

	// Create a new team member in database.
	err = db.Delete(&teamMember).Error
	if err != nil {
		// Return status 500 and database connection error.
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": true,
			"code":  repository.DATABASE_ERROR,
			"msg":   err.Error(),
		})
	}

	// Return status 201 and created team.
	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"error": false,
		"code":  repository.TEAM_MEMBER_REMOVED,
		"msg":   nil,
	})
}

// GetTeamMembers method to get all members from a team
// @Description Get all team members
// @Summary Get all team members
// @Tags Team
// @Accept json
// @Produce json
// @Param data body crud.TeamListMembers true "Team data"
// @Success 200 {string} status "ok"
// @Router /v1/team/members [get]
func GetTeamMembers(c *fiber.Ctx) error {
	listTeamMembers := &crud.TeamListMembers{}

	// Checking received data from JSON body.
	if err := c.BodyParser(listTeamMembers); err != nil {
		// Return status 400 and error message.
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": true,
			"code":  repository.INVALID_DATA,
			"msg":   err.Error(),
		})
	}

	// Create database connection.
	db := database.DB

	// Get user from database
	user := &models.User{}

	// Find user by ID
	err := db.Where(&models.User{ID: listTeamMembers.UserID}).First(&user).Error
	if err != nil {
		// Return status 500 and database connection error.
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": true,
			"code":  repository.USER_NOT_FOUND,
			"msg":   err.Error(),
		})
	}

	// Check if user is member of the team
	teamMember := &models.TeamHasMember{}

	err = db.Where(&models.TeamHasMember{TeamID: listTeamMembers.TeamID, MemberID: user.MemberID}).First(&teamMember).Error
	if err != nil {
		// Return status 500 and database connection error.
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": true,
			"code":  repository.DATABASE_ERROR,
			"msg":   err.Error(),
		})
	}

	if teamMember == nil {
		// Return status 500 and database connection error.
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": true,
			"code":  repository.TEAM_MEMBER_NOT_FOUND,
			"msg":   "User is not member of the team",
		})
	}

	teamMembers := []models.TeamHasMember{}

	// Get all members of the team
	err = db.Where(&models.TeamHasMember{TeamID: listTeamMembers.TeamID}).Find(&teamMembers).Error
	if err != nil {
		// Return status 500 and database connection error.
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": true,
			"code":  repository.DATABASE_ERROR,
			"msg":   err.Error(),
		})
	}

	// Return status 201 and created team.
	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"error": false,
		"code":  repository.TEAM_MEMBERS_OBTAINED,
		"msg":   nil,
	})
}
