package vault

import (
	"divSvault/utils"
	"encoding/json"
	"fmt"
	"log"
	"strings"
)

func GetSecret(vaultAddr, secretEngine, secretPath, token, namespace string) (string, error) {
	url := fmt.Sprintf("%s/v1/%s/data/%s", vaultAddr, secretEngine, secretPath)
	var payload *strings.Reader = nil
	data, err := utils.VaultRequest(url, "GET", token, namespace, payload)
	if err != nil {
		log.Printf("Error making req: %v", err)
	}
	secretData := data["data"].(map[string]interface{})["data"].(map[string]interface{})
	secret, err := json.MarshalIndent(secretData, "", "  ")
	if err != nil {
		return "", err
	}

	return string(secret), nil
}
