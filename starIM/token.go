package starIM

import (
	"github.com/dgrijalva/jwt-go"
	"time"
)

var jwtKey = []byte("02a119c540cdb7b9dafd2e9dd3d9f1eba6ec982bb1a549aaf87508d0086280b5")

type ClaimsACI struct {
	Ac string  //account
	P  string  //platform
	jwt.StandardClaims
}

func GenerateToken(ac string, p string) (string, error) {
	nowTime := time.Now()
	expireTime := nowTime.Add(time.Hour * 24 * 120)
	claims := &ClaimsACI{
		Ac: ac,
		P:  p,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: expireTime.Unix(),
			Issuer:    "Star",
		},
	}
	return jwt.NewWithClaims(jwt.SigningMethodHS384, claims).SignedString(jwtKey)
}

func ParseToken(token string) (*ClaimsACI, error) {
	tokenClaims, err := jwt.ParseWithClaims(token, &ClaimsACI{}, func(token *jwt.Token) (interface{}, error) {
		return jwtKey, nil
	})
	if err != nil {
		return nil, err
	}

	if tokenClaims != nil {
		if claims, ok := tokenClaims.Claims.(*ClaimsACI); ok && tokenClaims.Valid {
			return claims, nil
		}
	}

	return nil, err
}
