package main

import "GinChat/routes"

// @title           Swagger Example API
// @version         1.0
// @description     This is a Gin chat documentation
// @description     Base URL is in top
// @description     We_Don't_Know_What_Happened error usually is db error(access issue)
// @contact.name    API Support

// @license.name  ali.darzi.1354@gmail.com

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
