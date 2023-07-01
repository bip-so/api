package middlewares

import (
	"errors"
	"fmt"
	"gitlab.com/phonepost/bip-be-platform/internal/queries"
	"io/ioutil"
	"strconv"
	"strings"

	"github.com/dgrijalva/jwt-go"

	"github.com/gin-gonic/gin"

	"gitlab.com/phonepost/bip-be-platform/internal/models"
	"gitlab.com/phonepost/bip-be-platform/pkg/utils"
)

const bipTokenInValid = "bip-token-invalid"
const bipStudioRequestIdKey = "bip-studio-request-id"
const studioIdHeaderKey = "bip-studio-id"
const userIdKey = "currentUser"
const studioIdKey = "currentStudio"
const appBaseURL = "appBaseUrl"

/*
This Function will do the following .

1. Set "bip-studio-request-id" to header
2. Get and Set userID (UserID or None)
3. Get and Set studioID (studioID or 0)
*/
func SauronMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		setRequestId(c)
		getStudioIdFromHeaders(c)
		getUserFromHeader(c)
		setAppBaseURL(c)
		//RevisionMiddleware(c)
		c.Next()
	}
}

func setAppBaseURL(c *gin.Context) {
	scheme := "http"
	if c.Request.TLS != nil {
		scheme = "https"
	}
	url := scheme + "://" + c.Request.Host
	//fmt.Println("The Base URL: ", url)
	//fmt.Println(c.Request.URL.Path)
	c.Set(appBaseURL, url)
}

func RevisionMiddleware(c *gin.Context) {
	// Revision file contents will be only loaded once per process
	data, err := ioutil.ReadFile("REVISION")
	// If we cant read file, just skip to the next request handler
	// This is pretty much a NOOP middlware :)
	if err != nil {
		// Make sure to log error so it could be spotted
		fmt.Println("revision middleware error:", err)
	}
	/*
		res, err := ioutil.ReadFile("xxx")
		if err != nil {
		    switch err {
		    case os.ErrInvalid:
		        //Do stuff
		    case os.ErrPermission:
		        //Do stuff
		    case os.ErrNotExist:
		        //Do stuff
		    default:
		        fmt.Println(err)
		    }
		}https://stackoverflow.com/questions/23452157/how-do-i-check-for-specific-types-of-error-among-those-returned-by-ioutil-readfi?answertab=trending#tab-top
	*/

	// Clean up the value since it could contain line breaks
	revision := strings.TrimSpace(string(data))
	fmt.Println("revision middleware error:", revision)
	// Set out header value for each response
	c.Set("X-Revision", revision)
}

func setRequestId(c *gin.Context) {
	c.Header(bipStudioRequestIdKey, utils.NewNanoid())
}

func getStudioIdFromHeaders(c *gin.Context) {
	//"bip-studio-id"
	studioId := c.Request.Header.Get(studioIdHeaderKey)
	var studioIDFinal uint64
	if studioId != "" {
		studioIDFinal, _ = strconv.ParseUint(studioId, 10, 64)
		c.Set("currentStudio", studioIDFinal)
	}

	//ctx := context.Background()
	//ctx = context.WithValue(ctx, "currentStudioGo", "studioIDFinal dgsdfgs;ldfkjgsldfkjgs;lkd")
}

func getUserFromHeader(c *gin.Context) {
	auth := c.Request.Header.Get("Authorization")

	if auth == "" {
		c.Set("currentUser", nil)
		return
	}
	token := strings.TrimPrefix(auth, "Bearer ")
	if token == auth {
		c.Set("currentUser", nil)
		return
	}
	verifiedToken, err := VerifyToken(token)
	if err != nil {
		c.Set("currentUser", nil)
		return
	}

	// token validity
	err = TokenValidity(verifiedToken)
	fmt.Println(err)
	if err != nil {
		c.Set("currentUser", nil)
		return
	}
	var user *models.User
	claims, ok := verifiedToken.Claims.(jwt.MapClaims)

	if ok && verifiedToken.Valid {
		userId, _ := strconv.ParseUint(fmt.Sprintf("%.f", claims["userId"]), 10, 64)
		//fmt.Println("User id from token:", userId)
		//user, err = user2.App.Repo.GetUser(map[string]interface{}{"id": userId})
		user, err = queries.App.UserQueries.GetUserByID(userId)

		c.Set("currentUser", user)
		return
	}

	c.Set("currentUser", nil)
}

func GetUserFromLambdaHeader(header map[string]string) (*models.User, error) {
	auth := header["Authorization"]

	if auth == "" {
		return nil, errors.New("auth is empty")
	}
	token := strings.TrimPrefix(auth, "Bearer ")
	if token == auth {
		return nil, errors.New("token is equal to auth")
	}
	verifiedToken, err := VerifyToken(token)
	if err != nil {
		return nil, err
	}

	// token validity
	err = TokenValidity(verifiedToken)
	fmt.Println(err)
	if err != nil {
		return nil, err
	}
	var user *models.User
	claims, ok := verifiedToken.Claims.(jwt.MapClaims)

	if ok && verifiedToken.Valid {
		userId, _ := strconv.ParseUint(fmt.Sprintf("%.f", claims["userId"]), 10, 64)
		fmt.Println("User id from token:", userId)
		//user, err = user2.App.Repo.GetUser(map[string]interface{}{"id": userId})
		user, err = queries.App.UserQueries.GetUserByID(userId)

		if err != nil {
			return nil, err
		}
		return user, nil
	}

	return nil, errors.New("cannot get user")
}
