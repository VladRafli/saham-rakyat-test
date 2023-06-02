package helpers

import (
	"github.com/joho/godotenv"
	"github.com/labstack/echo/v4"
)

func LoadEnvironment(app *echo.Echo) {
	err := godotenv.Load()

	if err != nil {
		app.Logger.Fatal("Error loading .env file")
	}
}