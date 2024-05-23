package main

import (
	"fmt"

	"github.com/skye-tan/basket-manager/utils/api"
	"github.com/skye-tan/basket-manager/utils/auth"
	"github.com/skye-tan/basket-manager/utils/database"
)

func main() {
	fmt.Println("Generating secret-key...")
	auth.GenerateSecretKey()

	fmt.Println("Initializing database...")
	err := database.InitializeDatabase()
	if err != nil {
		panic(err)
	}

	fmt.Println("Initializing api...")
	api.InitializeApi()
}
