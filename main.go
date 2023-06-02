package main

import (
	"os"
	"sahamrakyat_test/database/migrations"
	"sahamrakyat_test/database/seeds"
	"sahamrakyat_test/helpers"
	"sahamrakyat_test/routes"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

func main() {
	app := echo.New()

	routes.Init(app)

	helpers.LoadEnvironment(app)

	db := helpers.ConnectDatabase(os.Getenv("DB_HOST"), os.Getenv("DB_USER"), os.Getenv("DB_PASS"), os.Getenv("DB_NAME"))
	migrations.Migrate(db)

	app.Use(middleware.CORS())
	// app.Use(middleware.CSRF()) // Not suitable for API used by mobile apps
	app.Use(middleware.Gzip())
	app.Use(middleware.Logger())
	app.Use(middleware.Recover())
	app.Use(middleware.RateLimiter(middleware.NewRateLimiterMemoryStore(20)))
	// app.Use(middleware.Secure()) // using X-Xss-Protection is known problematic (Find Chrome Bug report for about X-Xss-Protection)

	if os.Getenv("APP_ENV") == "development" {
		seeds.Seed(db)
		helpers.ClearCache()
	}

	app.Logger.Fatal(app.Start(":5000"))
}
