package controllers

import (
	"github.com/UIdealist/UIdealist-Members-Ms/app/crud"
	"github.com/UIdealist/UIdealist-Members-Ms/app/models"
	"github.com/UIdealist/UIdealist-Members-Ms/pkg/repository"
	"github.com/UIdealist/UIdealist-Members-Ms/platform/database"

	"github.com/gofiber/fiber/v2"
)

// CreateUser method to create a new user.
// @Description Create a new user given username and email
// @Summary Create a new user
// @Tags User
// @Accept json
// @Produce json
// @Param data body crud.UserCreate true "User data"
// @Success 201 {string} status "ok"
// @Router /v1/user [post]
func CreateUser(c *fiber.Ctx) error {
	userRequest := &crud.UserCreate{}

	// Checking received data from JSON body.
	if err := c.BodyParser(userRequest); err != nil {
		// Return status 400 and error message.
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": true,
			"code":  repository.INVALID_DATA,
			"msg":   err.Error(),
		})
	}

	// Create database connection.
	db := database.DB

	// Create a new user.
	user := &models.User{
		Username: userRequest.Username,
		Email:    userRequest.Email,
		Verified: false,
	}

	// Create a new user in database.
	err := db.Create(&user).Error
	if err != nil {
		// Return status 500 and database connection error.
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": true,
			"code":  repository.DATABASE_ERROR,
			"msg":   err.Error(),
		})
	}

	// Return status 201 and created user.
	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"error": false,
		"code":  repository.USER_CREATED,
		"msg":   nil,
		"user":  user,
	})
}

// CreateUserAnonymous method to create a new anonymous user.
// @Description Create a new user given username
// @Summary Create a new user
// @Tags User
// @Accept json
// @Produce json
// @Param data body crud.AnonymousUserCreate true "Anonymous User Data"
// @Success 201 {string} status "ok"
// @Router /v1/user/anonymous [post]
func CreateAnonymousUser(c *fiber.Ctx) error {
	anUserRequest := &crud.AnonymousUserCreate{}

	// Checking received data from JSON body.
	if err := c.BodyParser(anUserRequest); err != nil {
		// Return status 400 and error message.
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": true,
			"code":  repository.INVALID_DATA,
			"msg":   err.Error(),
		})
	}

	// Create database connection.
	db := database.DB

	// Create a new user.
	user := &models.AnonymousUser{
		TempName: anUserRequest.TempName,
	}

	// Create a new user in database.
	err := db.Create(&user).Error
	if err != nil {
		// Return status 500 and database connection error.
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": true,
			"code":  repository.DATABASE_ERROR,
			"msg":   err.Error(),
		})
	}

	// Return status 201 and created user.
	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"error": false,
		"code":  repository.USER_CREATED,
		"msg":   nil,
		"user":  user,
	})
}
