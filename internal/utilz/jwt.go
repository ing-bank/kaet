package utilz

import (
	"strings"

	"github.com/golang-jwt/jwt/v5"
)

type JWTClaims struct {
	Sub       string
	Namespace string
}

func GetJWTClaims(t string) *JWTClaims {
	jc := &JWTClaims{}

	if !LooksLikeJWT(t) {
		return jc
	}

	claims := jwt.MapClaims{}
	_, _ = jwt.ParseWithClaims(
		t,
		claims,
		nil,
	)
	if tokenSub, ok := claims["sub"]; ok {
		jc.Sub = tokenSub.(string)
	}

	if jc.Sub != "" {
		// Get namespace from Sub
		splitSub := strings.Split(jc.Sub, ":")
		if len(splitSub) > 3 {
			jc.Namespace = splitSub[2]
		}
	}

	if kClaims, exists := claims["kubernetes.io"]; exists && jc.Namespace == "" {
		jc.Namespace = kClaims.(map[string]any)["namespace"].(string)
	}

	return jc
}
