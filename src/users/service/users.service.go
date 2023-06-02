package service

import (
	"fmt"
	"net/http"
	"os"
	"strconv"
	"time"
	"sahamrakyat_test/database"
	"sahamrakyat_test/helpers"

	"github.com/go-playground/validator"
	"github.com/go-redis/cache/v8"
	"github.com/labstack/echo/v4"
)

func Create(c echo.Context) (*database.Users, *echo.HTTPError) {
	logger := helpers.InitLogger()

	logger.Info(helpers.ApacheFormatLogger(c.Request().Method, c.Request().URL.Host, c.Request().Host, c.RealIP(), c.Request().UserAgent(), time.Now().String()))

	user := &database.Users{}

	if err := c.Bind(user); err != nil {
		logger.Error(err)
		return nil, echo.NewHTTPError(http.StatusInternalServerError, "Failed to bind user.")
	}

	if err := validator.New().Struct(user); err != nil {
		logger.Error(err)

		validationErrors := []echo.Map{}

		if _, ok := err.(*validator.InvalidValidationError); ok {
			return nil, echo.NewHTTPError(http.StatusInternalServerError, "Invalid argument passed.")
		}

		for _, err := range err.(validator.ValidationErrors) {

			validationErrors = append(validationErrors, echo.Map{
				"namespace":       err.Namespace(),
				"field":           err.Field(),
				"structNamespace": err.StructNamespace(),
				"structField":     err.StructField(),
				"tag":             err.Tag(),
				"actualTag":       err.ActualTag(),
				"kind":            err.Kind(),
				"type":            err.Type(),
				"value":           err.Value(),
				"param":           err.Param(),
			})
		}

		return nil, echo.NewHTTPError(http.StatusBadRequest, echo.Map{
			"message": "Failed to validate user.",
			"error":   validationErrors,
		})
	}

	db := helpers.ConnectDatabase(os.Getenv("DB_HOST"), os.Getenv("DB_USER"), os.Getenv("DB_PASS"), os.Getenv("DB_NAME"))

	db.Create(user)

	logger.Info(fmt.Sprintf("[Response] %s %s %s %s %s", c.Response().Header().Get("Method"), c.Response().Header().Get("Host"), c.Response().Header().Get("RemoteAddr"), c.Response().Header().Get("UserAgent"), c.Response().Header().Get("Time")))

	return user, nil
}

func GetAll(c echo.Context) (*[]database.Users, *echo.HTTPError) {
	logger := helpers.InitLogger()

	logger.Info(helpers.ApacheFormatLogger(c.Request().Method, c.Request().URL.Host, c.Request().Host, c.RealIP(), c.Request().UserAgent(), time.Now().String()))

	var take int
	var skip int

	takeQuery := c.QueryParam("take")

	if takeQuery != "" {
		if _, err := strconv.ParseInt(takeQuery, 10, 32); err != nil {
			return nil, echo.NewHTTPError(http.StatusBadRequest, "Failed to parse take.")
		}
	}

	skipQuery := c.QueryParam("skip")

	if skipQuery != "" {
		if _, err := strconv.ParseInt(skipQuery, 10, 32); err != nil {
			return nil, echo.NewHTTPError(http.StatusBadRequest, "Failed to parse skip.")
		}
	}

	users := &[]database.Users{}

	db := helpers.ConnectDatabase(os.Getenv("DB_HOST"), os.Getenv("DB_USER"), os.Getenv("DB_PASS"), os.Getenv("DB_NAME"))

	if takeQuery != "" && skipQuery != "" {
		db.Limit(take).Offset(skip).Find(users)
	} else if takeQuery != "" {
		db.Limit(take).Find(users)
	} else if skipQuery != "" {
		db.Offset(skip).Find(users)
	} else {
		db.Find(users)
	}

	cacheClient := helpers.InitRedisCache()

	if err := cacheClient.Get(c.Request().Context(), "users", users); err == nil {
		return users, nil
	} else if err != nil {
		if err := cacheClient.Set(&cache.Item{
			Ctx:   c.Request().Context(),
			Key:   "users",
			Value: users,
		}); err != nil {
			fmt.Printf("Failed to cache users: %v", err)
		}
	}

	logger.Info(fmt.Sprintf("[Response] %s %s %s %s %s", c.Response().Header().Get("Method"), c.Response().Header().Get("Host"), c.Response().Header().Get("RemoteAddr"), c.Response().Header().Get("UserAgent"), c.Response().Header().Get("Time")))

	return users, nil
}

func Get(c echo.Context) (*database.Users, *echo.HTTPError) {
	logger := helpers.InitLogger()

	logger.Info(helpers.ApacheFormatLogger(c.Request().Method, c.Request().URL.Host, c.Request().Host, c.RealIP(), c.Request().UserAgent(), time.Now().String()))

	if c.Param("id") == "" {
		return nil, echo.NewHTTPError(http.StatusBadRequest, "User id is required.")
	}

	id, err := strconv.ParseUint(c.Param("id"), 10, 32)

	if err != nil {
		return nil, echo.NewHTTPError(http.StatusBadRequest, "Failed to parse user id.")
	}

	user := &database.Users{
		ID: uint(id),
	}

	if err := c.Bind(user); err != nil {
		return nil, echo.NewHTTPError(http.StatusInternalServerError, "Failed to bind user.")
	}

	db := helpers.ConnectDatabase(os.Getenv("DB_HOST"), os.Getenv("DB_USER"), os.Getenv("DB_PASS"), os.Getenv("DB_NAME"))

	if result := db.Find(user); result.Error != nil || result.RowsAffected < 1 {
		return nil, echo.NewHTTPError(http.StatusBadRequest, "User not found.")
	}

	cacheClient := helpers.InitRedisCache()

	if err := cacheClient.Get(c.Request().Context(), fmt.Sprintf("user:%d", id), user); err == nil {
		return user, nil
	} else if err != nil {
		if err := cacheClient.Set(&cache.Item{
			Ctx:   c.Request().Context(),
			Key:   fmt.Sprintf("user:%d", id),
			Value: user,
		}); err != nil {
			fmt.Printf("Failed to cache user: %v", err)
		}
	}

	logger.Info(fmt.Sprintf("[Response] %s %s %s %s %s", c.Response().Header().Get("Method"), c.Response().Header().Get("Host"), c.Response().Header().Get("RemoteAddr"), c.Response().Header().Get("UserAgent"), c.Response().Header().Get("Time")))

	return user, nil
}

func Update(c echo.Context) (*database.Users, *echo.HTTPError) {
	logger := helpers.InitLogger()

	logger.Info(helpers.ApacheFormatLogger(c.Request().Method, c.Request().URL.Host, c.Request().Host, c.RealIP(), c.Request().UserAgent(), time.Now().String()))

	if c.Param("id") == "" {
		return nil, echo.NewHTTPError(http.StatusBadRequest, "User id is required.")
	}

	id, err := strconv.ParseUint(c.Param("id"), 10, 32)

	if err != nil {
		return nil, echo.NewHTTPError(http.StatusBadRequest, "Failed to parse user id.")
	}

	user := &database.Users{
		ID: uint(id),
	}

	if err := c.Bind(&database.Users{}); err != nil {
		return nil, echo.NewHTTPError(http.StatusInternalServerError, "Failed to bind user.")
	}

	db := helpers.ConnectDatabase(os.Getenv("DB_HOST"), os.Getenv("DB_USER"), os.Getenv("DB_PASS"), os.Getenv("DB_NAME"))

	if result := db.First(&user); result.Error != nil || result.RowsAffected < 1 {
		return nil, echo.NewHTTPError(http.StatusBadRequest, "User not found.")
	}

	cacheClient := helpers.InitRedisCache()

	if err := cacheClient.Get(c.Request().Context(), fmt.Sprintf("user:%d", id), user); err == nil {
		cacheClient.Delete(c.Request().Context(), fmt.Sprintf("user:%d", id))
	}

	db.Save(user)

	logger.Info(fmt.Sprintf("[Response] %s %s %s %s %s", c.Response().Header().Get("Method"), c.Response().Header().Get("Host"), c.Response().Header().Get("RemoteAddr"), c.Response().Header().Get("UserAgent"), c.Response().Header().Get("Time")))

	return user, nil
}

func Delete(c echo.Context) (*database.Users, *echo.HTTPError) {
	logger := helpers.InitLogger()

	logger.Info(helpers.ApacheFormatLogger(c.Request().Method, c.Request().URL.Host, c.Request().Host, c.RealIP(), c.Request().UserAgent(), time.Now().String()))

	if c.Param("id") == "" {
		return nil, echo.NewHTTPError(http.StatusBadRequest, "User id is required.")
	}

	id, err := strconv.ParseUint(c.Param("id"), 10, 32)

	if err != nil {
		return nil, echo.NewHTTPError(http.StatusBadRequest, "Failed to parse user id.")
	}

	user := &database.Users{
		ID: uint(id),
	}

	if err := c.Bind(&database.Users{}); err != nil {
		return nil, echo.NewHTTPError(http.StatusInternalServerError, "Failed to bind user.")
	}

	db := helpers.ConnectDatabase(os.Getenv("DB_HOST"), os.Getenv("DB_USER"), os.Getenv("DB_PASS"), os.Getenv("DB_NAME"))

	if result := db.First(&user); result.Error != nil || result.RowsAffected < 1 {
		return nil, echo.NewHTTPError(http.StatusBadRequest, "User not found.")
	}

	cacheClient := helpers.InitRedisCache()

	if err := cacheClient.Get(c.Request().Context(), fmt.Sprintf("user:%d", id), user); err == nil {
		cacheClient.Delete(c.Request().Context(), fmt.Sprintf("user:%d", id))
	}

	db.Delete(user)

	logger.Info(fmt.Sprintf("[Response] %s %s %s %s %s", c.Response().Header().Get("Method"), c.Response().Header().Get("Host"), c.Response().Header().Get("RemoteAddr"), c.Response().Header().Get("UserAgent"), c.Response().Header().Get("Time")))

	return user, nil
}
