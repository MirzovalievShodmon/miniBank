package main

import (
	"fmt"

	"github.com/MirzovalievShodmon/miniBank.git/internal/controller"
	"github.com/MirzovalievShodmon/miniBank.git/internal/db"
)

// The Power of Justice
func main() {
	if err := db.InitConnection(); err != nil {
		fmt.Println("Error during database connection initialization: ", err.Error())
		return
	}

	if err := db.RunMigrations(); err != nil {
		fmt.Println("Error during database migrations: ", err.Error())
		return
	}

	controller.InitRoutes()

	if err := db.CloseConnection(); err != nil {
		fmt.Println("Error during database connection close: ", err.Error())
		return
	}
}
