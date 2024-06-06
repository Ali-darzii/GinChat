package main

import "GinChat/routes"

func main() {
	router := routes.Urls()

	err := router.Run(":8080")
	if err != nil {
		panic("Can't start the server")
	}

}
