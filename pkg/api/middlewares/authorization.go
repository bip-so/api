package middlewares

import (
	"fmt"
	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	"gitlab.com/phonepost/bip-be-platform/internal/models"
	"gitlab.com/phonepost/bip-be-platform/internal/queries"
	"gitlab.com/phonepost/bip-be-platform/pkg/configs"
	"gitlab.com/phonepost/bip-be-platform/pkg/stores/redis"
	"net/http"
	"strconv"
	"strings"
)

func TokenAuthorizationMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		auth := c.Request.Header.Get("Authorization")

		if auth == "" {
			c.String(http.StatusUnauthorized, "No Authorization header provided")
			c.Abort()
			return
		}
		token := strings.TrimPrefix(auth, "Bearer ")
		if token == auth {
			c.String(http.StatusUnauthorized, "No Bearer Token Found Could not find bearer token in Authorization header")
			c.Abort()
			return
		}

		// Adding check for "invalidated tokens"
		// token-sessions is the key
		rc := redis.RedisClient()
		ctx := redis.GetBgContext()
		val, err := rc.Get(ctx, "token-sessions:"+token).Result()
		if val == "false" {
			if err.Error() == "Token is expired" {
				c.Header(bipTokenInValid, "TOKEN_INVALID")
			}
			// This custom Talk to GM and CC
			c.String(406, "verifiedToken: Token was Deleted by User.")
			c.Abort()
			return
		}

		//token verify
		verifiedToken, err := VerifyToken(token)
		if err != nil {
			if err.Error() == "Token is expired" {
				c.Header(bipTokenInValid, "TOKEN_INVALID")
			}
			c.String(http.StatusUnauthorized, "verifiedToken: Token is not verified.")
			c.Abort()
			return
		}

		// token validity
		err = TokenValidity(verifiedToken)
		fmt.Println(err)
		if err != nil {
			fmt.Println(err.Error())
			c.String(http.StatusUnauthorized, "TokenValidity: Token is not verified.")
			c.Abort()
			return
		}
		var user *models.User
		claims, ok := verifiedToken.Claims.(jwt.MapClaims)

		if ok && verifiedToken.Valid {
			userId, _ := strconv.ParseUint(fmt.Sprintf("%.f", claims["userId"]), 10, 64)
			//user, err = user2.App.Repo.GetUser(map[string]interface{}{"id": userId})
			user, err = queries.App.UserQueries.GetUserByID(userId)
			c.Set("currentUser", user)
		}
		c.Next()
	}
}

func VerifyToken(token string) (*jwt.Token, error) {
	verified_token, err := jwt.Parse(token, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("Unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(configs.GetConfigString("ACCESS_SECRET")), nil
	})
	if err != nil {
		return nil, err
	}
	return verified_token, nil
}

func TokenValidity(token *jwt.Token) error {
	var err error
	if err != nil {
		return err
	}
	if _, ok := token.Claims.(jwt.Claims); !ok && !token.Valid {
		return err
	}
	return nil
}
