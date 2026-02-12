package lib

import (
	"errors"
	"strings"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt"
	"github.com/google/uuid"
	"github.com/spf13/viper"
)

type ReqHeaderAuth struct {
	Authorization string `reqHeader:"authorization"`
}

var Claim Claims

type Claims struct {
	jwt.StandardClaims
	UserID      *uuid.UUID `json:"user_id"`
	Role        string     `json:"role"`
	Permissions []string   `json:"permissions"`
}

// ClaimsJWT func
func ClaimsJWT(accesToken *string) (jwt.MapClaims, error) {
	token, _, err := new(jwt.Parser).ParseUnverified(*accesToken, jwt.MapClaims{})
	if err != nil {
		return nil, err
	}
	claims, _ := token.Claims.(jwt.MapClaims)

	timeNow := time.Now().Unix()
	timeSessions := int64(claims["exp"].(float64))
	if timeSessions < timeNow {
		return claims, err
	}

	return claims, nil
}

func ParseJwt(cookie string, secretKey ...string) (*Claims, error) {
	secret := viper.GetString("JWT_SECRET")
	if len(secretKey) > 0 {
		secret = secretKey[0]
	}

	claims := Claims{}
	if token, err := jwt.ParseWithClaims(cookie, &claims, func(token *jwt.Token) (interface{}, error) {
		return []byte(secret), nil
	}); err != nil || !token.Valid {
		return nil, err
	}

	return &claims, nil
}

func GetToken(c *fiber.Ctx) string {
	token := ""

	bearerToken := new(ReqHeaderAuth)
	if err := c.ReqHeaderParser(bearerToken); err == nil {
		token, _ = bearerToken.GetBearerToken()
	}

	tokenArr := []string{
		"access_token",
		"token",
	}

	if len(token) == 0 {
		for _, v := range tokenArr {
			if c.Cookies(v) != "" {
				token = c.Cookies(v)
				break
			}
		}
	}

	return token
}

func (h *ReqHeaderAuth) GetBearerToken() (string, error) {
	if len(h.Authorization) == 0 {
		err := errors.New("authorization header not found")
		return "", err
	}

	authorization := strings.Split(h.Authorization, " ")
	if strings.ToLower(authorization[0]) != "bearer" {
		err := errors.New("not a bearer token")
		return "", err
	}

	return authorization[1], nil
}

func GenerateAccessToken(userID *uuid.UUID) (string, error) {
	claims := &Claims{
		UserID: userID,
		StandardClaims: jwt.StandardClaims{
			Issuer:    StringUUID(userID),
			ExpiresAt: time.Now().Add(time.Hour * time.Duration(viper.GetInt("LOGIN_SESSION"))).Unix(),
		},
	}

	tokens := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	return tokens.SignedString([]byte(viper.GetString("JWT_SECRET")))
}

func GenerateRefreshToken(userID *uuid.UUID) (string, error) {
	claims := &Claims{
		UserID: userID,
		StandardClaims: jwt.StandardClaims{
			Issuer:    StringUUID(userID),
			ExpiresAt: time.Now().Add(time.Hour * time.Duration(viper.GetInt("REFRESH_SESSION"))).Unix(),
		},
	}

	tokens := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	return tokens.SignedString([]byte(viper.GetString("JWT_SECRET")))
}

func GetXUserID(c *fiber.Ctx) *uuid.UUID {
	claims := Claims{}

	if token, err := jwt.ParseWithClaims(GetToken(c), &claims, func(token *jwt.Token) (interface{}, error) {
		return []byte(viper.GetString("SECRET_KEY")), nil
	}); err == nil && token.Valid {
		return StrToUUID(claims.Issuer)
	}

	return nil
}

func GetXRole(c *fiber.Ctx) string {
	claims := Claims{}

	if token, err := jwt.ParseWithClaims(GetToken(c), &claims, func(token *jwt.Token) (interface{}, error) {
		return []byte(viper.GetString("JWT_SECRET")), nil
	}); err == nil && token.Valid {
		return claims.Role
	}

	return ""
}

func GetXPermission(c *fiber.Ctx) []string {
	claims := Claims{}

	if token, err := jwt.ParseWithClaims(GetToken(c), &claims, func(token *jwt.Token) (interface{}, error) {
		return []byte(viper.GetString("JWT_SECRET")), nil
	}); err == nil && token.Valid {
		return claims.Permissions
	}

	return nil
}
