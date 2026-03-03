package jwt

import "github.com/golang-jwt/jwt/v5"

type Claims struct {
	Id        int64
	Username  string
	Email     string
	Status    string
	Roles     []string
	IsDeleted bool
	jwt.RegisteredClaims
}
