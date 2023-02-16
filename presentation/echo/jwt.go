package echo

import (
	"errors"
	"strings"
	"time"

	"github.com/golang-jwt/jwt"
	"github.com/labstack/echo/v4"
)

type jwtCustumClaims struct {
	Username string `json:"username"`
	jwt.StandardClaims
}

type GenerateTokenParam struct {
	Username string
}

func GenerateToken(p GenerateTokenParam) (token string, err error) {
	claims := &jwtCustumClaims{
		p.Username,
		jwt.StandardClaims{
			ExpiresAt: time.Now().Add(time.Hour * 72).Unix(),
			IssuedAt:  time.Now().Unix(),
			Issuer:    *jwtIssuer,
		},
	}

	newToken := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return newToken.SignedString([]byte(*jwtSecret))
}

func verifyToken(token *jwt.Token) (claims *jwtCustumClaims, err error) {
	claims = token.Claims.(*jwtCustumClaims)

	if !claims.VerifyIssuer(*jwtIssuer, true) {
		err = errors.New("token is invalid")
		return
	}

	if !claims.VerifyExpiresAt(time.Now().Unix(), true) {
		err = errors.New("token is expired")
		return
	}

	return
}

func CheckJWT(c echo.Context) (claims *jwtCustumClaims, err error) {
	authHeader := c.Request().Header.Get("Authorization")
	if authHeader != "" && strings.HasPrefix(authHeader, "Bearer ") {
		tokenString := strings.TrimPrefix(authHeader, "Bearer ")

		var token *jwt.Token
		token, err = jwt.ParseWithClaims(
			tokenString,
			&jwtCustumClaims{},
			func(token *jwt.Token) (interface{}, error) {
				if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
					return nil, errors.New("unexpected signing method")
				}
				return []byte(*jwtSecret), nil
			},
		)
		if err != nil {
			return nil, err
		}
		return verifyToken(token)
	}
	return nil, errors.New("bearer token not specified")
}
