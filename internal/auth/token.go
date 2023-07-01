package auth

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/dgrijalva/jwt-go"

	"gitlab.com/phonepost/bip-be-platform/pkg/configs"
	"gitlab.com/phonepost/bip-be-platform/pkg/stores/redis"
	"gitlab.com/phonepost/bip-be-platform/pkg/utils"
)

type TokenDetails struct {
	AccessToken    string `json:"accessToken"`
	RefreshToken   string `json:"refreshToken"`
	AccessTokenID  string `json:"accessTokenID"`
	RefreshTokenID string `json:"refreshTokenID"`
}

// Returns The Access Token & The Refresh Token
// Access Token : 24 hours
// Refresh Token: 7 days
func CreateToken(userid uint64) (*TokenDetails, error) {
	tokens := &TokenDetails{}
	tokens.AccessTokenID = utils.GetUUID()
	tokens.RefreshTokenID = utils.GetUUID()
	var err error

	// Access Token
	accessTokenClaim := jwt.MapClaims{}
	accessTokenClaim["authorized"] = true

	accessTokenClaim["accessUUID"] = tokens.AccessTokenID
	accessTokenClaim["userId"] = userid
	accessTokenClaim["exp"] = time.Now().Add(time.Hour * 24 * 7).Unix()
	//accessTokenClaim["exp"] = time.Now().Add(time.Second * 60).Unix()
	accessToken := jwt.NewWithClaims(jwt.SigningMethodHS256, accessTokenClaim)
	tokens.AccessToken, err = accessToken.SignedString([]byte(configs.GetConfigString("ACCESS_SECRET")))
	if err != nil {
		return nil, err
	}
	// Refresh Token
	refreshTokentClaims := jwt.MapClaims{}
	refreshTokentClaims["refreshUUID"] = tokens.RefreshTokenID
	refreshTokentClaims["userId"] = userid
	refreshTokentClaims["exp"] = time.Now().Add(time.Hour * 24 * 15).Unix()
	//refreshTokentClaims["exp"] = time.Now().Add(time.Second * 300).Unix()
	rt := jwt.NewWithClaims(jwt.SigningMethodHS256, refreshTokentClaims)
	tokens.RefreshToken, err = rt.SignedString([]byte(configs.GetConfigString("REFRESH_SECRET")))
	if err != nil {
		return nil, err
	}

	// save to redis
	rc := redis.RedisClient()
	ctx := redis.GetBgContext()

	tokenWithDataud := make(map[string]interface{})
	tokenWithDataud["refreshToken"] = tokens.RefreshToken
	tokenWithDataud["accessToken"] = tokens.AccessToken
	tokenWithDataud["accessTokenId"] = tokens.AccessTokenID
	tokenWithDataud["refreshTokenId"] = tokens.AccessTokenID
	tokenWithDataud["userId"] = userid

	tokenDetailsJson, err := json.Marshal(tokens)
	if err != nil {
		fmt.Println(err)
	}
	_ = rc.Set(ctx, "authtokens:"+tokens.AccessTokenID, tokenDetailsJson, time.Hour*24*7).Err()
	// Setting ACCESS TOKEN ID IN NEW KEY
	_ = rc.Set(ctx, "token-sessions:"+tokens.AccessToken, "true", time.Hour*24*7*30).Err()
	//_ = rc.Set(ctx, "authtokens:"+strconv.Itoa(int(userid)), tokenDetailsJson, 0).Err()
	return tokens, nil
}

func RefreshToken(token string) (string, error) {
	return "", nil
}

// Gets a String Token and Returns JJwt token/
func ValidTokenFromBearerToken(token string) (*jwt.Token, error) {
	verified_token, err := jwt.Parse(token, func(token *jwt.Token) (interface{}, error) {
		//Make sure that the token method conform to "SigningMethodHMAC"
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

func CheckTokenValidity(token *jwt.Token) error {
	var err error
	if err != nil {
		return err
	}
	if _, ok := token.Claims.(jwt.Claims); !ok && !token.Valid {
		return err
	}
	return nil
}

// Deleted
//func GetUserFromJWTToken(token jwt.Token) (*models.User, error) {
//	var thisuser *models.User
//	claims, ok := token.Claims.(jwt.MapClaims)
//	if ok && token.Valid {
//		var uerr error
//		userId, _ := strconv.ParseUint(fmt.Sprintf("%.f", claims["userId"]), 10, 64)
//
//		user, uerr := user2.App.Repo.GetUser(map[string]interface{}{"id": userId})
//		if uerr != nil {
//			return nil, uerr
//		}
//		thisuser = user
//	}
//	return thisuser, nil
//}

func ParseJWTToken(token string) (result map[string]interface{}, valid bool, err error) {

	jwttoken, err := jwt.Parse(token, func(token *jwt.Token) (interface{}, error) {
		claims, ok := token.Claims.(jwt.MapClaims)
		println("claims:", claims)
		if !ok {
			return false, nil
		}
		result = claims
		return token, nil
	})

	var isValid bool
	if jwttoken != nil && jwttoken.Valid {
		isValid = true
	} else {
		isValid = false
	}
	return result, isValid, err
}

func CreateCustomJWTToken(userID uint64, payload map[string]interface{}) (string, error) {
	tokens := &TokenDetails{}
	tokens.AccessTokenID = utils.GetUUID()

	var err error

	// Access Token
	accessTokenClaim := jwt.MapClaims{}
	accessTokenClaim["authorized"] = true

	accessTokenClaim["accessUUID"] = tokens.AccessTokenID
	accessTokenClaim["userId"] = userID
	for k, v := range payload {
		accessTokenClaim[k] = v
	}
	accessTokenClaim["exp"] = time.Now().Add(time.Hour * 24 * 7).Unix()
	accessToken := jwt.NewWithClaims(jwt.SigningMethodHS256, accessTokenClaim)
	tokens.AccessToken, err = accessToken.SignedString([]byte(configs.GetConfigString("ACCESS_SECRET")))
	if err != nil {
		return "", err
	}

	return tokens.AccessToken, nil
}
