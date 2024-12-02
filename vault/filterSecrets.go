package vault

import (
	"strings"
)

func FilterSecrets(secrets []string, query string) []string {
	var filtered []string
	for _, secret := range secrets {
		if strings.Contains(secret, query) {
			filtered = append(filtered, secret)
		}
	}
	return filtered
}
