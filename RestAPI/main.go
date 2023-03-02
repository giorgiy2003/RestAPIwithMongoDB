package main

import (
	Handler "myapp/internal/handlers"
	Logic "myapp/internal/logic"
	Repository "myapp/internal/repository"

	"github.com/labstack/echo"
)

func main() {
	Logic.InitLogger()
	router := echo.New()
	router.Use(Handler.ConnectDB)
	router.GET("/person", Handler.GetPersons)
	router.GET("/person/:id", Handler.GetById)
	router.POST("/person", Handler.PostPerson)
	router.DELETE("/person/:id", Handler.DeleteById)
	router.PUT("/person/:id", Handler.UpdatePersonById)
	router.Logger.Fatal(router.Start(":8080"))
	defer Repository.Disconnect()
}
