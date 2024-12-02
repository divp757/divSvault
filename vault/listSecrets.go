package vault

import (
	"divSvault/utils"
	"fmt"
	"log"
	"strings"
)

func ListSecrets(vaultAddr, secretEngine, token, apiPath, namespace string) ([]string, error) {
	url := fmt.Sprintf("%s/v1/%s/metadata/%s", vaultAddr, secretEngine, apiPath)
	var payload *strings.Reader = nil
	data, err := utils.VaultRequest(url, "LIST", token, namespace, payload)
	if err != nil {
		log.Printf("Error making req: %v", err)
	}
	keys := data["data"].(map[string]interface{})["keys"].([]interface{})
	secrets := make([]string, len(keys))
	for i, key := range keys {
		secrets[i] = key.(string)
	}

	return secrets, nil
}
