package controllers

import (
	"idealist/app/crud"
	"idealist/app/models"
	"idealist/pkg/repository"
	"idealist/platform/database"

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
// @Summary Create a new team
// @Tags Team
// @Accept json
// @Produce json
// @Param data body crud.TeamAddMember true "New member data"
// @Success 201 {string} status "ok"
// @Router /v1/team/members/add [post]
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
		"code":  repository.TEAM_CREATED,
		"msg":   nil,
	})
}
