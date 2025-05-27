package kubernetes

import (
	"github.com/golang-jwt/jwt/v5"
	"github.com/projectdiscovery/gologger"
)

func GetNamespaceFromSAToken(sat string) string {
	namespace := "default"

	if sat == "" {
		return namespace
	}

	claims := jwt.MapClaims{}
	_, err := jwt.ParseWithClaims(
		sat,
		claims,
		func(t *jwt.Token) (interface{}, error) {
			return []byte("secret"), nil
		},
	)
	if err != nil {
		gologger.Debug().
			Str("error", err.Error()).
			Msg("could not parse service account token\n")
	}

	if tokenNs := claims["kubernetes.io"].(map[string]any)["namespace"]; tokenNs != "" {
		namespace = tokenNs.(string)
	}

	return namespace
}
