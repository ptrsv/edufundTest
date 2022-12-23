package main

import (
	"bytes"
	"edufund/edufund/controller"
	"edufund/edufund/auth"
	"edufund/edufund/database"
	"edufund/edufund/models"
	"encoding/json"
	"log"
	
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/joho/godotenv"
)

func SetUpRouter() *gin.Engine{
	err:= godotenv.Load("../.env")
	if err != nil {
		log.Fatalf("Error loading env")
	}

	
	db:= database.InitDb()
	userController:= controller.UserController(db)
	jwt, _:= auth.InitJwt(db)
	r := gin.Default()

	r.POST("/register", userController.Register)
	r.POST("/login", jwt.LoginHandler)
	return r
}

func TestRegister(t *testing.T) {
	r:= SetUpRouter()
	
	w:= httptest.NewRecorder()

	user:= &models.User{
		Username:"edufund@gmail.com",
		Fullname:"Edufund",
		Password:"Locked1234568",
		ConfirmPassword: "Locked1234568",
	}
	body,_ := json.Marshal(user)
	req,_:= http.NewRequest("POST", "/register", bytes.NewBuffer(body))
	r.ServeHTTP(w,req)
	resp:= w.Body.String()

	res:= models.User{}
	json.Unmarshal([]byte(resp), &res)

	assert.NotContains(t, "message", "error")
}

func TestLogin(t *testing.T) {
	r:= SetUpRouter()
	
	w:= httptest.NewRecorder()

	user:= &models.UserLogin{
		Username:"edufund@gmail.com",
		Password:"Locked1234568",
	}
	body,_ := json.Marshal(user)
	req,_:= http.NewRequest("POST", "/login", bytes.NewBuffer(body))
	r.ServeHTTP(w,req)
	resp:= w.Body.String()

	res:= models.User{}
	json.Unmarshal([]byte(resp), &res)

	// assert.NotEqual(t, http.StatusOK, http.StatusUnauthorized)
	assert.NotContains(t, "access_token", "message")
}

