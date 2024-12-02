package vault

import (
	"divSvault/utils"
	"encoding/json"
	"fmt"
	"strings"
)

func WriteSecret(vaultAddr, secretEngine, secretPath string, secretData map[string]interface{}, token, namespace string) error {
	url := fmt.Sprintf("%s/v1/%s/data/%s", vaultAddr, secretEngine, secretPath)
	data := map[string]interface{}{
		"data": secretData,
	}
	jsonData, err := json.Marshal(data)
	if err != nil {
		return err
	}
	utils.VaultRequest(url, "POST", token, namespace, strings.NewReader(string(jsonData)))

	return nil
}
