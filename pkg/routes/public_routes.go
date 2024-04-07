package routes

import (
	"idealist/app/controllers"

	"github.com/gofiber/fiber/v2"
)

// PublicRoutes func for describe group of public routes.
func PublicRoutes(a *fiber.App) {
	// Create routes group.
	route := a.Group("/api/v1")

	// Routes for POST method:
	route.Post("/user", controllers.CreateUser)
	route.Post("/user/anonymous", controllers.CreateAnonymousUser)
	route.Post("/team", controllers.CreateTeam)
	route.Post("/team/members", controllers.AddTeamMember)
	route.Get("/team/members", controllers.GetTeamMembers)
	route.Delete("/team/members", controllers.RemoveTeamMember)
}
