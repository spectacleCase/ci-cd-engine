package request

import (
	jwt "github.com/golang-jwt/jwt/v5"
)

// CustomClaims structure
type CustomClaims struct {
	BaseClaims
	jwt.RegisteredClaims
}

type BaseClaims struct {
	ID uint
}
