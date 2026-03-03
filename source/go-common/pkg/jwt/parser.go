package jwt

import "github.com/golang-jwt/jwt/v5"

func ParseToken(tokenStr string, key []byte) (*Claims, error) {
	claims := &Claims{}
	token, err := jwt.ParseWithClaims(tokenStr, claims, func(token *jwt.Token) (any, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			var zero *Claims
			return zero, ErrInvalidToken
		}
		return key, nil
	})

	if err != nil || !token.Valid {
		var zero *Claims
		return zero, ErrInvalidToken
	}

	return claims, nil
}
