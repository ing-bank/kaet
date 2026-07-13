package utilz

import (
	"encoding/base64"
	"strings"
)

func LooksLikeJWT(i string) bool {
	parts := strings.Split(i, ".")
	if len(parts) < 2 {
		return false
	}

	for i, p := range parts {
		// last index of a jwt
		if i != 2 && p == "" {
			return false
		}
		_, err := base64.RawURLEncoding.DecodeString(p)
		if err != nil {
			return false
		}
	}

	return true
}
