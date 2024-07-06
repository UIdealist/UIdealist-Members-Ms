package main

import (
	"os"

	"github.com/UIdealist/UIdealist-Members-Ms/pkg/configs"
	"github.com/UIdealist/UIdealist-Members-Ms/pkg/middleware"
	"github.com/UIdealist/UIdealist-Members-Ms/pkg/routes"
	"github.com/UIdealist/UIdealist-Members-Ms/pkg/utils"
	"github.com/UIdealist/UIdealist-Members-Ms/platform/database"

	"github.com/gofiber/fiber/v2"

	_ "github.com/UIdealist/UIdealist-Members-Ms/docs" // load API Docs files (Swagger)

	_ "github.com/joho/godotenv/autoload" // load .env file automatically

	accessconnector "github.com/UIdealist/Uidealist-Access-Ms/connector" // load access microservice connector
)

// @UIdealist API
// @version 1.0
// @description UIdealist Member project API Docs.
// @termsOfService http://swagger.io/terms/
// @contact.name API Support
// @contact.email edgardanielgd123@gmail.com
// @license.name Apache 2.0
// @license.url http://www.apache.org/licenses/LICENSE-2.0.html
// @BasePath /api
func main() {
	// Define Fiber config.
	config := configs.FiberConfig()

	// Create database connection.
	database.OpenDBConnection()

	// Define a new Fiber app with config.
	app := fiber.New(config)

	// Initialize access microservice connector.
	accessconnector.Init()

	// Middlewares.
	middleware.FiberMiddleware(app) // Register Fiber's middleware for app.

	// Routes.
	routes.SwaggerRoute(app)  // Register a route for API Docs (Swagger).
	routes.PublicRoutes(app)  // Register a public routes for app.
	routes.PrivateRoutes(app) // Register a private routes for app.
	routes.NotFoundRoute(app) // Register route for 404 Error.

	// Start server (with or without graceful shutdown).
	if os.Getenv("STAGE_STATUS") == "dev" {
		utils.StartServer(app)
	} else {
		utils.StartServerWithGracefulShutdown(app)
	}

}
