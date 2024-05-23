package api

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"github.com/labstack/echo/v4"
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
		return c.String(http.StatusInternalServerError, "Expected application/json content-type.")
	}

	// Extracting request's json boby.
	content := make(map[string]interface{})
	err := json.NewDecoder(c.Request().Body).Decode(&content)
	if err != nil {
		return c.String(http.StatusBadRequest, "Invalid json body was provided.")
	}

	// Extracting the data from reqeust's json body.
	data, ok := content["data"].(string)
	if !ok {
		return c.String(http.StatusBadRequest, "No state was provided.")
	}

	// Extracting the state from reqeust's json body.
	state, ok := content["state"].(string)
	if !ok {
		return c.String(http.StatusBadRequest, "No state was provided.")
	}

	// Checking for validity of the provided state.
	if state != "COMPLETED" && state != "PENDING" {
		return c.String(http.StatusBadRequest, "Invalid state was provided.")
	}

	// Creating the new record in the database.
	id := database.CreateBasket(data, state)

	return c.String(http.StatusOK, fmt.Sprintf("Basket with id[%d] was created successfully.", id))
}

// PATCH "/basket/:id"
func updateBasket(c echo.Context) error {
	fmt.Println("Running PATCH /basket/:id")

	// Extracting the id parameter.
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		return c.String(http.StatusBadRequest, "Invalid id was provided.")
	}

	// Checking request's content type.
	content_type := c.Request().Header[echo.HeaderContentType]
	if !contains(content_type, "application/json") {
		return c.String(http.StatusInternalServerError, "Expected application/json content-type.")
	}

	// Extracting request's json boby.
	content := make(map[string]interface{})
	err = json.NewDecoder(c.Request().Body).Decode(&content)
	if err != nil {
		return c.String(http.StatusBadRequest, "Invalid json body was provided.")
	}

	// Extracting the data from reqeust's json body.
	data, ok := content["data"].(string)
	if !ok {
		return c.String(http.StatusBadRequest, "No state was provided.")
	}

	// Extracting the state from reqeust's json body.
	state, ok := content["state"].(string)
	if !ok {
		return c.String(http.StatusBadRequest, "No state was provided.")
	}

	// Checking for validity of the provided state.
	if state != "COMPLETED" && state != "PENDING" {
		return c.String(http.StatusBadRequest, "Invalid state was provided.")
	}

	// Updating the record in the database.
	database.UpdateBasket(uint(id), data, state)

	return c.String(http.StatusOK, "Basket was updateted successfully.")
}

// GET "/basket/:id"
func getBasket(c echo.Context) error {
	fmt.Println("Running GET /basket/:id")

	// Extracting the id parameter.
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		return c.String(http.StatusBadRequest, "Invalid id was provided.")
	}

	// Getting basket with the provided id.
	basket, err := database.GetBasket(uint(id))
	if err != nil {
		return c.String(http.StatusBadRequest, "Invalid id was provided.")
	}

	return c.JSON(http.StatusOK, basket)
}

// DELETE "/basket/:id"
func deleteBasket(c echo.Context) error {
	fmt.Println("Running DELETE /basket/:id.")

	// Extracting the id parameter.
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		return c.String(http.StatusBadRequest, "Invalid id was provided.")
	}

	// Deleting the record from the database.
	database.DeleteBasket(uint(id))

	return c.String(http.StatusOK, "The basket was deleted successfully.")
}
