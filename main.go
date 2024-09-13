package main

import (
	"GinChat/routes"
)

// @title           Swagger Example API
// @version         1.0
// @description     This is a Gin chat documentation

// @contact.name   API Support

// @license.name  Apache 2.0

// @host      localhost:8080
// @BasePath  /api/v1

// @securityDefinitions.basic  Json Web Token (jwt)

// @externalDocs.description  OpenAPI
func main() {
	router := routes.Urls()

	err := router.Run(":8080")
	if err != nil {
		panic("Can't start the server")
	}

}
