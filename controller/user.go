package controller

import (
	"fmt"
	"edufund/edufund/models"
	"net/http"

	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type UserRepo struct {
	Db * gorm.DB
}

func UserController(db *gorm.DB) *UserRepo {
	db.AutoMigrate(&models.User{})
	return &UserRepo{Db: db}
}

func (repo *UserRepo) Register(c *gin.Context) {
	var user models.User

	if c.BindJSON(&user) == nil {
		var existing models.User
		models.GetUser(repo.Db, &existing, user.Username)
		if(existing.ID > 0) {
			c.JSON(http.StatusBadRequest, gin.H{"error": "username already taken"})
			return
		}

		validate:= models.Validate(&user)
		if validate != nil {
			c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": validate.Error()})
			return
		}

		passwordBytes, _ := bcrypt.GenerateFromPassword([]byte(user.Password), 10)
		user.Password = string(passwordBytes)

		err:=  models.Register(repo.Db, &user)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err})
			return
		}
		c.JSON(http.StatusOK, gin.H{"message":"user successfuly registered"})
	} else {
		fmt.Println(user)
		c.JSON(http.StatusBadRequest, user)
	}
}