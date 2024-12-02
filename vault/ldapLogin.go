package vault

import (
	"bytes"
	"divSvault/utils"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/manifoldco/promptui"
)

func LdapLogin(vaultConfig utils.VaultConfig) (string, error) {
	usernamePrompt := promptui.Prompt{
		Label:   "Username",
		Default: vaultConfig.Username,
	}
	username, err := usernamePrompt.Run()
	if err != nil {
		return "", err
	}

	passwordPrompt := promptui.Prompt{
		Label: "Password",
		Mask:  '*',
	}
	password, err := passwordPrompt.Run()
	if err != nil {
		return "", err
	}

	loginData := map[string]string{
		"password": password,
	}

	loginDataBytes, err := json.Marshal(loginData)
	if err != nil {
		return "", err
	}

	url := fmt.Sprintf("%s/v1/auth/ldap/login/%s", vaultConfig.Address, username)
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(loginDataBytes))
	if err != nil {
		return "", err
	}
	req.Header.Set("Content-Type", "application/json")

	if vaultConfig.Namespace != "" {
		req.Header.Add("X-Vault-Namespace", vaultConfig.Namespace)
	}

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("failed to login: %s", resp.Status)
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	var data map[string]interface{}
	if err := json.Unmarshal(body, &data); err != nil {
		return "", err
	}

	return data["auth"].(map[string]interface{})["client_token"].(string), nil
}
