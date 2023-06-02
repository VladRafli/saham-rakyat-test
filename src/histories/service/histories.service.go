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
	"gorm.io/gorm/clause"
)

func Create(c echo.Context) (*database.Histories, *echo.HTTPError) {
	logger := helpers.InitLogger()

	logger.Info(helpers.ApacheFormatLogger(c.Request().Method, c.Request().URL.Host, c.Request().Host, c.RealIP(), c.Request().UserAgent(), time.Now().String()))

	history := &database.Histories{}

	if err := c.Bind(history); err != nil {
		fmt.Println(err)
		return nil, echo.NewHTTPError(http.StatusInternalServerError, "Failed to bind history.")
	}

	if err := validator.New().Struct(history); err != nil {
		fmt.Println(err)

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
			"message": "Failed to validate history.",
			"error":   validationErrors,
		})
	}

	db := helpers.ConnectDatabase(os.Getenv("DB_HOST"), os.Getenv("DB_USER"), os.Getenv("DB_PASS"), os.Getenv("DB_NAME"))

	db.Create(history)

	logger.Info(fmt.Sprintf("[Response] %s %s %s %s %s", c.Response().Header().Get("Method"), c.Response().Header().Get("Host"), c.Response().Header().Get("RemoteAddr"), c.Response().Header().Get("UserAgent"), c.Response().Header().Get("Time")))

	return history, nil
}

func GetAll(c echo.Context) (*[]database.Histories, *echo.HTTPError) {
	logger := helpers.InitLogger()

	logger.Info(helpers.ApacheFormatLogger(c.Request().Method, c.Request().URL.Host, c.Request().Host, c.RealIP(), c.Request().UserAgent(), time.Now().String()))

	var take int
	var skip int

	takeQuery := c.QueryParam("take")

	if takeQuery != "" {
		if val, err := strconv.ParseInt(takeQuery, 10, 32); err != nil {
			return nil, echo.NewHTTPError(http.StatusBadRequest, "Failed to parse take.")
		} else {
			take = int(val)
		}
	} else {
		take = -1
	}

	skipQuery := c.QueryParam("skip")

	if skipQuery != "" {
		if val, err := strconv.ParseInt(skipQuery, 10, 32); err != nil {
			return nil, echo.NewHTTPError(http.StatusBadRequest, "Failed to parse skip.")
		} else {
			skip = int(val)
		}
	} else {
		skip = -1
	}
	
	histories := &[]database.Histories{}

	db := helpers.ConnectDatabase(os.Getenv("DB_HOST"), os.Getenv("DB_USER"), os.Getenv("DB_PASS"), os.Getenv("DB_NAME"))

	db.Limit(take).Offset(skip).Preload(clause.Associations).Find(histories)

	cacheClient := helpers.InitRedisCache()

	if err := cacheClient.Get(c.Request().Context(), "histories", histories); err == nil {
		return histories, nil
	} else if err != nil {
		if err := cacheClient.Set(&cache.Item{
			Ctx:   c.Request().Context(),
			Key:   "histories",
			Value: histories,
		}); err != nil {
			fmt.Printf("Failed to cache histories: %v\n", err)
		}
	}

	logger.Info(fmt.Sprintf("[Response] %s %s %s %s %s", c.Response().Header().Get("Method"), c.Response().Header().Get("Host"), c.Response().Header().Get("RemoteAddr"), c.Response().Header().Get("UserAgent"), c.Response().Header().Get("Time")))

	return histories, nil
}

func Get(c echo.Context) (*database.Histories, *echo.HTTPError) {
	logger := helpers.InitLogger()

	logger.Info(helpers.ApacheFormatLogger(c.Request().Method, c.Request().URL.Host, c.Request().Host, c.RealIP(), c.Request().UserAgent(), time.Now().String()))

	if c.Param("id") == "" {
		return nil, echo.NewHTTPError(http.StatusBadRequest, "History id is required.")
	}
	
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)

	if err != nil {
		return nil, echo.NewHTTPError(http.StatusBadRequest, "Failed to parse history id.")
	}

	history := &database.Histories{
		ID: uint(id),
	}

	if err := c.Bind(history); err != nil {
		return nil, echo.NewHTTPError(http.StatusInternalServerError, "Failed to bind history.")
	}

	db := helpers.ConnectDatabase(os.Getenv("DB_HOST"), os.Getenv("DB_USER"), os.Getenv("DB_PASS"), os.Getenv("DB_NAME"))

	if result := db.Find(history); result.Error != nil || result.RowsAffected < 0 {
		return nil, echo.NewHTTPError(http.StatusNotFound, "History not found.")
	}

	cacheClient := helpers.InitRedisCache()

	if err := cacheClient.Get(c.Request().Context(), fmt.Sprintf("history:%d", id), history); err == nil {
		return history, nil
	} else if err != nil {
		if err := cacheClient.Set(&cache.Item{
			Ctx:   c.Request().Context(),
			Key:   fmt.Sprintf("history:%d", id),
			Value: history,
		}); err != nil {
			fmt.Printf("Failed to cache history: %v\n", err)
		}
	}

	logger.Info(fmt.Sprintf("[Response] %s %s %s %s %s", c.Response().Header().Get("Method"), c.Response().Header().Get("Host"), c.Response().Header().Get("RemoteAddr"), c.Response().Header().Get("UserAgent"), c.Response().Header().Get("Time")))
	
	return history, nil
}

func Update(c echo.Context) (*database.Histories, *echo.HTTPError) {
	logger := helpers.InitLogger()

	logger.Info(helpers.ApacheFormatLogger(c.Request().Method, c.Request().URL.Host, c.Request().Host, c.RealIP(), c.Request().UserAgent(), time.Now().String()))

	if c.Param("id") == "" {
		return nil, echo.NewHTTPError(http.StatusBadRequest, "History id is required.")
	}

	id, err := strconv.ParseUint(c.Param("id"), 10, 32)

	if err != nil {
		return nil, echo.NewHTTPError(http.StatusBadRequest, "Failed to parse history id.")
	}

	history := &database.Histories{
		ID: uint(id),
	}

	if err := c.Bind(history); err != nil {
		fmt.Println(err)
		return nil, echo.NewHTTPError(http.StatusInternalServerError, "Failed to bind history.")
	}

	db := helpers.ConnectDatabase(os.Getenv("DB_HOST"), os.Getenv("DB_USER"), os.Getenv("DB_PASS"), os.Getenv("DB_NAME"))

	if result := db.Find(history); result.Error != nil || result.RowsAffected < 0 {
		return nil, echo.NewHTTPError(http.StatusNotFound, "History not found.")
	}

	cacheClient := helpers.InitRedisCache()

	if err := cacheClient.Get(c.Request().Context(), fmt.Sprintf("history:%d", id), history); err == nil {
		cacheClient.Delete(c.Request().Context(), fmt.Sprintf("history:%d", id))
	}

	db.Save(history)

	logger.Info(fmt.Sprintf("[Response] %s %s %s %s %s", c.Response().Header().Get("Method"), c.Response().Header().Get("Host"), c.Response().Header().Get("RemoteAddr"), c.Response().Header().Get("UserAgent"), c.Response().Header().Get("Time")))

	return history, nil
}

func Delete(c echo.Context) (*database.Histories, *echo.HTTPError) {
	logger := helpers.InitLogger()

	logger.Info(helpers.ApacheFormatLogger(c.Request().Method, c.Request().URL.Host, c.Request().Host, c.RealIP(), c.Request().UserAgent(), time.Now().String()))

	if c.Param("id") == "" {
		return nil, echo.NewHTTPError(http.StatusBadRequest, "History id is required.")
	}

	id, err := strconv.ParseUint(c.Param("id"), 10, 32)

	if err != nil {
		return nil, echo.NewHTTPError(http.StatusBadRequest, "Failed to parse history id.")
	}

	history := &database.Histories{
		ID: uint(id),
	}

	if err := c.Bind(history); err != nil {
		fmt.Println(err)
		return nil, echo.NewHTTPError(http.StatusInternalServerError, "Failed to bind history.")
	}

	db := helpers.ConnectDatabase(os.Getenv("DB_HOST"), os.Getenv("DB_USER"), os.Getenv("DB_PASS"), os.Getenv("DB_NAME"))

	if result := db.Find(history); result.Error != nil || result.RowsAffected < 0 {
		return nil, echo.NewHTTPError(http.StatusNotFound, "History not found.")
	}

	cacheClient := helpers.InitRedisCache()

	if err := cacheClient.Get(c.Request().Context(), fmt.Sprintf("history:%d", id), history); err == nil {
		cacheClient.Delete(c.Request().Context(), fmt.Sprintf("history:%d", id))
	}

	db.Delete(history)

	logger.Info(fmt.Sprintf("[Response] %s %s %s %s %s", c.Response().Header().Get("Method"), c.Response().Header().Get("Host"), c.Response().Header().Get("RemoteAddr"), c.Response().Header().Get("UserAgent"), c.Response().Header().Get("Time")))

	return history, nil
}
