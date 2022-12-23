package auth

import (
	"edufund/edufund/models"
	"os"
	"time"

	jwt "github.com/appleboy/gin-jwt/v2"
	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

var identityKey = "id"
var identityName = "name"

func InitJwt(db *gorm.DB) (*jwt.GinJWTMiddleware, error) {
	JWT_KEY:= os.Getenv("JWT_KEY")
	JWT_REALM:= os.Getenv("JWT_REALM")

	db.AutoMigrate(&models.User{})

	authMiddleware, err:= jwt.New(&jwt.GinJWTMiddleware{
		Realm: JWT_REALM,
		Key: []byte(JWT_KEY),
		Timeout: 1*24*time.Hour,
		MaxRefresh: 1*24*time.Hour,
		IdentityKey: identityKey,
		PayloadFunc: func(data interface{}) jwt.MapClaims {
			if v, ok:= data.(*models.UserPayload); ok {
				return jwt.MapClaims{
					identityKey: v.Username,
					identityName: v.Fullname,
				}
			}
			return jwt.MapClaims{}
		},
		IdentityHandler: func(c *gin.Context) interface{} {
			claims:= jwt.ExtractClaims(c)
			return &models.UserPayload{
				Username: claims[identityKey].(string),
				Fullname: claims[identityName].(string),
			}
		},
		Authenticator: func(c *gin.Context) (interface{}, error) {
			var userLogin models.UserLogin 
			if err:= c.ShouldBind(&userLogin); err != nil {
				return "", jwt.ErrMissingLoginValues
			}
			username:= userLogin.Username
			password:= userLogin.Password

			var user models.User

			ValidateEmail:= models.ValidateEmaiil(username)
			if ValidateEmail != nil {
				return nil, ValidateEmail
			}

			validatePassword:= models.ValidatePassword(password, "")
			if validatePassword != nil {
				return nil, validatePassword
			}


			err:= models.GetUser(db, &user, username)
			if err == nil {
				err:=  bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))
				if err == nil {
					return &models.UserPayload{
						Username: username,
						Fullname: user.Fullname,
					}, nil
				}
			}
			return nil, jwt.ErrFailedAuthentication
		},
		Authorizator: func(data interface{}, c *gin.Context) bool {
			if _, ok:= data.(*models.UserPayload); ok {
				return true
			}
			return false
		},
		LoginResponse: func (c *gin.Context, code int, token string, expire time.Time)  {
			c.JSON(code, gin.H{
				"access_token": token,
				"expires_in": expire.Format(time.RFC3339),
				"message": "login success",
			})
		},
	})
	return authMiddleware, err
}