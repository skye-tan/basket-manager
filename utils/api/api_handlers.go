package api

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"github.com/labstack/echo/v4"
	"github.com/skye-tan/basket-manager/utils"
	"github.com/skye-tan/basket-manager/utils/database"
)

func InitializeApi() {
	e := echo.New()

	e.GET("/basket", getBaskets)

	e.POST("/basket", createBasket)

	e.PATCH("/basket/:id", updateBasket)

	e.GET("/basket/:id", getBasket)

	e.DELETE("/basket/:id", deleteBasket)

	e.Start("0.0.0.0:8081")
}

func contains(slice []string, element string) bool {
	for counter := range slice {
		if slice[counter] == element {
			return true
		}
	}
	return false
}

// GET "/basket"
func getBaskets(c echo.Context) error {
	fmt.Println("Running GET /basket")

	// Getting all of the baskets.
	baskets := database.GetBaskets()

	return c.JSON(http.StatusOK, baskets)
}

// POST "/basket"
func createBasket(c echo.Context) error {
	fmt.Println("Running POST /basket")

	// Checking request's content type.
	content_type := c.Request().Header[echo.HeaderContentType]
	if !contains(content_type, "application/json") {
		return c.String(http.StatusInternalServerError, custom_error.INVALID_CONTENT_TYPE)
	}

	// Extracting request's json boby.
	content := make(map[string]interface{})
	err := json.NewDecoder(c.Request().Body).Decode(&content)
	if err != nil {
		return c.String(http.StatusBadRequest, custom_error.INVALID_BODY_FORMAT)
	}

	// Extracting the data from reqeust's json body.
	data, ok := content["data"].(string)
	if !ok {
		return c.String(http.StatusBadRequest, custom_error.MISSING_DATA)
	}

	// Extracting the state from reqeust's json body.
	state, ok := content["state"].(string)
	if !ok {
		return c.String(http.StatusBadRequest, custom_error.MISSING_STATE)
	}

	// Creating the new record in the database.
	id, err := database.CreateBasket(data, state)
	if err != nil {
		return c.String(http.StatusOK, err.Error())
	}

	return c.String(http.StatusOK, fmt.Sprintf("created successfully with id = %d", id))
}

// PATCH "/basket/:id"
func updateBasket(c echo.Context) error {
	fmt.Println("Running PATCH /basket/:id")

	// Extracting the id parameter.
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		return c.String(http.StatusBadRequest, custom_error.INVALID_ID)
	}

	// Checking request's content type.
	content_type := c.Request().Header[echo.HeaderContentType]
	if !contains(content_type, "application/json") {
		return c.String(http.StatusInternalServerError, custom_error.INVALID_CONTENT_TYPE)
	}

	// Extracting request's json boby.
	content := make(map[string]interface{})
	err = json.NewDecoder(c.Request().Body).Decode(&content)
	if err != nil {
		return c.String(http.StatusBadRequest, custom_error.INVALID_ID)
	}

	// Extracting the data from reqeust's json body.
	data, ok := content["data"].(string)
	if !ok {
		return c.String(http.StatusBadRequest, custom_error.MISSING_DATA)
	}

	// Extracting the state from reqeust's json body.
	state, ok := content["state"].(string)
	if !ok {
		return c.String(http.StatusBadRequest, custom_error.MISSING_STATE)
	}

	// Updating the record in the database.
	err = database.UpdateBasket(uint(id), data, state)
	if err != nil {
		return c.String(http.StatusBadRequest, err.Error())
	}

	return c.String(http.StatusOK, "updateted successfully")
}

// GET "/basket/:id"
func getBasket(c echo.Context) error {
	fmt.Println("Running GET /basket/:id")

	// Extracting the id parameter.
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		return c.String(http.StatusBadRequest, custom_error.INVALID_ID)
	}

	// Getting basket with the provided id.
	basket, err := database.GetBasket(uint(id))
	if err != nil {
		return c.String(http.StatusBadRequest, err.Error())
	}

	return c.JSON(http.StatusOK, basket)
}

// DELETE "/basket/:id"
func deleteBasket(c echo.Context) error {
	fmt.Println("Running DELETE /basket/:id.")

	// Extracting the id parameter.
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		return c.String(http.StatusBadRequest, custom_error.INVALID_ID)
	}

	// Deleting the record from the database.
	err = database.DeleteBasket(uint(id))
	if err != nil {
		return c.String(http.StatusBadRequest, err.Error())
	}

	return c.String(http.StatusOK, "deleted successfully")
}
