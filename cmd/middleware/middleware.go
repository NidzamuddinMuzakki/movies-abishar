package middleware

import (
	"net/http"
	"strings"

	"github.com/NidzamuddinMuzakki/movies-abishar/go-lib-common/cache"
	common "github.com/NidzamuddinMuzakki/movies-abishar/go-lib-common/registry"

	"github.com/NidzamuddinMuzakki/movies-abishar/common/response"
	"github.com/NidzamuddinMuzakki/movies-abishar/common/util"
	"github.com/NidzamuddinMuzakki/movies-abishar/config"
	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
)

type IMiddleware interface {
	AuthJWT() gin.HandlerFunc
}

type middleware struct {
	common common.IRegistry
}

func NewMiddleware(common common.IRegistry) IMiddleware {
	return &middleware{
		common: common,
	}
}

func (m *middleware) AuthJWT() gin.HandlerFunc {
	return func(c *gin.Context) {

		ctx := c.Request.Context()

		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.JSON(http.StatusUnauthorized, response.Unauthorised(ctx))
			c.AbortWithStatus(http.StatusUnauthorized)
			return
		}

		authHeaderParts := strings.Split(authHeader, " ")
		if len(authHeaderParts) != 2 || authHeaderParts[0] != "Bearer" {
			c.JSON(http.StatusUnauthorized, response.Unauthorised(ctx))
			c.AbortWithStatus(http.StatusUnauthorized)
			return
		}

		tokenString := authHeaderParts[1]

		isValid, err := validateJWT(tokenString)
		if err != nil {
			c.JSON(http.StatusUnauthorized, response.Unauthorised(ctx))
			c.AbortWithStatus(http.StatusUnauthorized)
			return
		}

		if !isValid {
			c.JSON(http.StatusUnauthorized, response.Unauthorised(ctx))
			c.AbortWithStatus(http.StatusUnauthorized)
			return
		}
		tokenData := util.ReadDataToken(tokenString)
		// fmt.Println(tokenData.Uuid, "nidzam")
		var TOkenData interface{}
		err = m.common.GetCache().Get(ctx, cache.Key(tokenData.Uuid), &TOkenData)
		// fmt.Println(err, "nidzam")
		if err != nil {
			// fmt.Println(err, "errr NIdzam")
			if err.Error() != "redis: nil" {
				c.JSON(http.StatusInternalServerError, "please try again")
				c.AbortWithStatus(http.StatusInternalServerError)
				return
			}
		} else {
			c.JSON(http.StatusUnauthorized, response.Unauthorised(ctx))
			c.AbortWithStatus(http.StatusUnauthorized)
			return
		}
		c.Next()
	}
}

func validateJWT(tokenString string) (bool, error) {
	var mySigningKey = []byte(config.Cold.JwtSecretKey)

	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		return mySigningKey, nil
	})

	if err != nil {
		return false, err
	}

	if _, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		return true, nil
	} else {
		return false, nil
	}
}
