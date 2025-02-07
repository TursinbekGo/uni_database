package helper

import (
	"errors"
	"strings"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/spf13/cast"
)

type TokenInfo struct {
	UserID     string `json:"user_id"`
	ClientType string `json:"client_type"`
	PlatformID string `json:"platform_id"`
}
type TadqiqotchiTokenInfo struct {
	TadqiqotchiID string `json:"tadqiqotchi_id"`
}
type OtherTokenInfo struct {
	OtherID string `json:"tadqiqotchi_id"`
}
type OqituvchiTokenInfo struct {
	OqituvchiID string `json:"oqituvchi_id"`
}

// GenerateJWT ...
func GenerateJWT(m map[string]interface{}, tokenExpireTime time.Duration, tokenSecretKey string) (tokenString string, err error) {
	var token *jwt.Token

	token = jwt.New(jwt.SigningMethodHS256)

	claims := token.Claims.(jwt.MapClaims)

	for key, value := range m {
		claims[key] = value
	}

	claims["iat"] = time.Now().Unix()
	claims["exp"] = time.Now().Add(tokenExpireTime).Unix()

	tokenString, err = token.SignedString([]byte(tokenSecretKey))
	if err != nil {
		return "", err
	}

	return tokenString, nil
}

func ParseClaimsForOther(token string, secretKey string) (result OtherTokenInfo, err error) {
	var claims jwt.MapClaims

	claims, err = ExtractClaims(token, secretKey)
	if err != nil {
		return result, err
	}

	result.OtherID = cast.ToString(claims["user_id"])
	if len(result.OtherID) <= 0 {
		err = errors.New("cannot parse 'user_id' field")
		return result, err
	}

	return
}
func ParseClaimsForOqituvchi(token string, secretKey string) (result OqituvchiTokenInfo, err error) {
	var claims jwt.MapClaims

	claims, err = ExtractClaims(token, secretKey)
	if err != nil {
		return result, err
	}

	result.OqituvchiID = cast.ToString(claims["user_id"])
	if len(result.OqituvchiID) <= 0 {
		err = errors.New("cannot parse 'user_id' field")
		return result, err
	}

	return
}
func ParseClaimsForTadqiqotchi(token string, secretKey string) (result TadqiqotchiTokenInfo, err error) {
	var claims jwt.MapClaims

	claims, err = ExtractClaims(token, secretKey)
	if err != nil {
		return result, err
	}

	result.TadqiqotchiID = cast.ToString(claims["user_id"])
	if len(result.TadqiqotchiID) <= 0 {
		err = errors.New("cannot parse 'user_id' field")
		return result, err
	}

	return
}

// ExtractClaims extracts claims from given token
func ExtractClaims(tokenString string, tokenSecretKey string) (jwt.MapClaims, error) {
	var (
		token *jwt.Token
		err   error
	)

	token, err = jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		// check token signing method etc
		return []byte(tokenSecretKey), nil
	})

	if err != nil {
		return nil, err
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !(ok && token.Valid) {
		return nil, errors.New("invalid token")
	}

	return claims, nil
}

// ExtractToken checks and returns token part of input string
func ExtractToken(bearer string) (token string, err error) {
	strArr := strings.Split(bearer, " ")
	if len(strArr) == 2 {
		return strArr[1], nil
	}
	return token, errors.New("wrong token format")
}

func ParseClaims(token string, secretKey string) (result TokenInfo, err error) {
	var claims jwt.MapClaims

	claims, err = ExtractClaims(token, secretKey)
	if err != nil {
		return result, err
	}

	result.UserID = cast.ToString(claims["user_id"])
	result.ClientType = cast.ToString(claims["client_type"])
	result.PlatformID = cast.ToString(claims["platform_id"])
	if len(result.UserID) <= 0 {
		err = errors.New("cannot parse 'user_id' field")
		return result, err
	}

	return
}
