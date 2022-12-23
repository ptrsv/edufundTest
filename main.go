package main

import (
	"edufund/edufund/auth"
	"edufund/edufund/controller"
	"edufund/edufund/database"
	"fmt"
	"log"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func main() {
	err:= godotenv.Load(".env")
	if err != nil {
		log.Fatalf("Error loading env")
	}

	r:= gin.Default()

	db:= database.InitDb()
	userController:= controller.UserController(db)
	jwt, _:= auth.InitJwt(db)

	r.POST("/register", userController.Register)
	r.POST("/login", jwt.LoginHandler)

	r.Run()
	fmt.Println("server running baby")
}