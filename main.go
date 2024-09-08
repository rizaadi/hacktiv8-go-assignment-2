package main

import (
	"hacktiv8-go-assignment-2/database"
	"hacktiv8-go-assignment-2/routers"
)

func main() {
	database.StartDB()

	var PORT = ":8080"
	routers.StartServer().Run(PORT)
}
