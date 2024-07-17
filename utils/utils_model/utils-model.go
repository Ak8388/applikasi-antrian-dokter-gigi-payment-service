package utilsmodel

import "github.com/golang-jwt/jwt/v5"

type JWTClaims struct {
	jwt.RegisteredClaims
	Email string
	Id    string
	Role  string
	Name  string
}
