package api

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"github.com/labstack/echo/v4"
	custom_error "github.com/skye-tan/basket-manager/utils"
	"github.com/skye-tan/basket-manager/utils/auth"
	"github.com/skye-tan/basket-manager/utils/database"
)

func InitializeApi() {
	e := echo.New()

	e.POST("/register", register)

	e.POST("/login", login)

	e.GET("/basket", getBaskets)

	e.POST("/basket", createBasket)

	e.PATCH("/basket/:id", updateBasket)

	e.GET("/basket/:id", getBasket)

	e.DELETE("/basket/:id", deleteBasket)

	e.Start("0.0.0.0:8081")
}

func extractToken(c echo.Context) string {
	authorization := c.Request().Header.Get(echo.HeaderAuthorization)

	token_string := strings.TrimPrefix(authorization, "sso-jwt ")

	return token_string
}

// POST "/register"
func register(c echo.Context) error {
	fmt.Println("Running POST /register")

	// Checking request's content type.
	content_type := c.Request().Header.Get(echo.HeaderContentType)
	if content_type != "application/json" {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"message": custom_error.INVALID_CONTENT_TYPE,
		})
	}

	// Extracting request's json boby.
	content := make(map[string]interface{})
	err := json.NewDecoder(c.Request().Body).Decode(&content)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"message": custom_error.INVALID_BODY_FORMAT,
		})
	}

	// Extracting the username from reqeust's json body.
	username, ok := content["username"].(string)
	if !ok {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"message": custom_error.MISSING_USERNAME,
		})
	}

	// Extracting the password from reqeust's json body.
	password, ok := content["password"].(string)
	if !ok {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"message": custom_error.MISSING_PASSWORD,
		})
	}

	// Adding the new user in the database.
	err = database.AddUser(username, password)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"message": err.Error(),
		})
	}

	return c.JSON(http.StatusCreated, map[string]string{
		"message": "registered successfully",
	})
}

// POST "/login"
func login(c echo.Context) error {
	fmt.Println("Running POST /login")

	// Checking request's content type.
	content_type := c.Request().Header.Get(echo.HeaderContentType)
	if content_type != "application/json" {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"message": custom_error.INVALID_CONTENT_TYPE,
		})
	}

	// Extracting request's json boby.
	content := make(map[string]interface{})
	err := json.NewDecoder(c.Request().Body).Decode(&content)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"message": custom_error.INVALID_BODY_FORMAT,
		})
	}

	// Extracting the username from reqeust's json body.
	username, ok := content["username"].(string)
	if !ok {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"message": custom_error.MISSING_USERNAME,
		})
	}

	// Extracting the password from reqeust's json body.
	password, ok := content["password"].(string)
	if !ok {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"message": custom_error.MISSING_PASSWORD,
		})
	}

	// Checking for the validity of the credentials.
	user, err := database.GetUser(username)
	if err != nil || user.Password != password {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"message": err.Error(),
		})
	}

	// Generating jwt.
	token_string, err := auth.CreateToken(user.ID)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{
			"message": custom_error.TOKEN_FAILURE,
		})
	}

	return c.JSON(http.StatusOK, map[string]string{
		"message": "logged in successfully",
		"token":   token_string,
	})
}

// GET "/basket"
func getBaskets(c echo.Context) error {
	fmt.Println("Running GET /basket")

	// Validating jwt.
	user_id, err := auth.VerifyToken(extractToken(c))
	if err != nil {
		return c.JSON(http.StatusUnauthorized, map[string]string{
			"message": err.Error(),
		})
	}

	// Getting all of the baskets.
	baskets := database.GetBaskets(user_id)

	return c.JSON(http.StatusOK, baskets)
}

// POST "/basket"
func createBasket(c echo.Context) error {
	fmt.Println("Running POST /basket")

	// Validating jwt.
	user_id, err := auth.VerifyToken(extractToken(c))
	if err != nil {
		return c.JSON(http.StatusUnauthorized, map[string]string{
			"message": err.Error(),
		})
	}

	// Checking request's content type.
	content_type := c.Request().Header.Get(echo.HeaderContentType)
	if content_type != "application/json" {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"message": custom_error.INVALID_CONTENT_TYPE,
		})
	}

	// Extracting request's json boby.
	content := make(map[string]interface{})
	err = json.NewDecoder(c.Request().Body).Decode(&content)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"message": custom_error.INVALID_BODY_FORMAT,
		})
	}

	// Extracting the data from reqeust's json body.
	data, ok := content["data"].(map[string]interface{})
	if !ok {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"message": custom_error.MISSING_DATA,
		})
	}

	// Extracting the state from reqeust's json body.
	state, ok := content["state"].(string)
	if !ok {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"message": custom_error.MISSING_STATE,
		})
	}

	// Creating the new record in the database.
	id, err := database.CreateBasket(user_id, data, state)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"message": err.Error(),
		})
	}

	return c.JSON(http.StatusCreated, map[string]string{
		"message":   "created successfully",
		"basket_id": fmt.Sprint(id),
	})
}

// PATCH "/basket/:id"
func updateBasket(c echo.Context) error {
	fmt.Println("Running PATCH /basket/:id")

	// Validating jwt.
	user_id, err := auth.VerifyToken(extractToken(c))
	if err != nil {
		return c.JSON(http.StatusUnauthorized, map[string]string{
			"message": err.Error(),
		})
	}

	// Extracting the id parameter.
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"message": custom_error.INVALID_ID,
		})
	}

	// Checking request's content type.
	content_type := c.Request().Header.Get(echo.HeaderContentType)
	if content_type != "application/json" {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"message": custom_error.INVALID_CONTENT_TYPE,
		})
	}

	// Extracting request's json boby.
	content := make(map[string]interface{})
	err = json.NewDecoder(c.Request().Body).Decode(&content)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"message": custom_error.INVALID_ID,
		})
	}

	// Extracting the data from reqeust's json body.
	data, ok := content["data"].(map[string]interface{})
	if !ok {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"message": custom_error.MISSING_DATA,
		})
	}

	// Extracting the state from reqeust's json body.
	state, ok := content["state"].(string)
	if !ok {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"message": custom_error.MISSING_STATE,
		})
	}

	// Updating the record in the database.
	err = database.UpdateBasket(user_id, uint(id), data, state)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"message": err.Error(),
		})
	}

	return c.JSON(http.StatusOK, map[string]string{
		"message": "updateted successfully",
	})
}

// GET "/basket/:id"
func getBasket(c echo.Context) error {
	fmt.Println("Running GET /basket/:id")

	// Validating jwt.
	user_id, err := auth.VerifyToken(extractToken(c))
	if err != nil {
		return c.JSON(http.StatusUnauthorized, map[string]string{
			"message": err.Error(),
		})
	}

	// Extracting the id parameter.
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"message": custom_error.INVALID_ID,
		})
	}

	// Getting basket with the provided id.
	basket, err := database.GetBasket(user_id, uint(id))
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"message": err.Error(),
		})
	}

	return c.JSON(http.StatusOK, basket)
}

// DELETE "/basket/:id"
func deleteBasket(c echo.Context) error {
	fmt.Println("Running DELETE /basket/:id.")

	// Validating jwt.
	user_id, err := auth.VerifyToken(extractToken(c))
	if err != nil {
		return c.JSON(http.StatusUnauthorized, map[string]string{
			"message": err.Error(),
		})
	}

	// Extracting the id parameter.
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"message": custom_error.INVALID_ID,
		})
	}

	// Deleting the record from the database.
	err = database.DeleteBasket(user_id, uint(id))
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"message": err.Error(),
		})
	}

	return c.JSON(http.StatusOK, map[string]string{
		"message": "deleted successfully",
	})
}
