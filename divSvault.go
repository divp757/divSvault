package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/manifoldco/promptui"
)

func main() {
	// Define flags for the CLI
	vaultAddr := flag.String("vault-addr", getEnv("VAULT_ADDR", "http://127.0.0.1:8200"), "Vault server address")
	secretEngine := flag.String("secret-engine", "secret", "Secret engine in Vault")
	token := flag.String("token", getEnv("VAULT_TOKEN", ""), "Vault token")
	namespace := flag.String("namespace", getEnv("VAULT_NAMESPACE", ""), "Vault namespace")
	outputFormat := flag.String("output-format", "json", "Output format (json or text)")

	flag.Parse()

	if *token == "" {
		log.Fatal("Vault token is required")
	}

	for {
		// Create a prompt to select an action
		actionPrompt := promptui.Select{
			Label: "Select Action",
			Items: []string{"List Secrets", "Add Secret", "Update Secret"},
		}

		_, action, err := actionPrompt.Run()
		if err != nil {
			log.Printf("Prompt failed %v\n", err)
			os.Exit(0)
		}

		switch action {
		case "List Secrets":
			listAndDisplaySecrets(*vaultAddr, *secretEngine, *token, *namespace, *outputFormat)
		case "Add Secret":
			addSecret(*vaultAddr, *secretEngine, *token, *namespace)
		case "Update Secret":
			updateSecret(*vaultAddr, *secretEngine, *token, *namespace)
		}
	}
}

func getEnv(key, fallback string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return fallback
}

func listAndDisplaySecrets(vaultAddr, secretEngine, token, namespace, outputFormat string) {
	secrets, err := listSecrets(vaultAddr, secretEngine, token, namespace)
	if err != nil {
		log.Printf("Error listing secrets: %v", err)
		return
	}

	searchPrompt := promptui.Prompt{
		Label: "Search Secret",
	}

	searchQuery, err := searchPrompt.Run()
	if err != nil {
		log.Printf("Prompt failed %v\n", err)
		return
	}

	filteredSecrets := filterSecrets(secrets, searchQuery)

	selectPrompt := promptui.Select{
		Label: "Select Secret",
		Items: filteredSecrets,
	}

	_, result, err := selectPrompt.Run()
	if err != nil {
		log.Printf("Prompt failed %v\n", err)
		return
	}

	secret, err := getSecret(vaultAddr, secretEngine, result, token, namespace)
	if err != nil {
		log.Printf("Error getting secret: %v", err)
		return
	}

	displaySecret(secret, outputFormat)
}

func listSecrets(vaultAddr, secretEngine, token, namespace string) ([]string, error) {
	url := fmt.Sprintf("%s/v1/%s/metadata", vaultAddr, secretEngine)
	req, err := http.NewRequest("LIST", url, nil)
	if err != nil {
		return nil, err
	}

	req.Header.Add("X-Vault-Token", token)
	if namespace != "" {
		req.Header.Add("X-Vault-Namespace", namespace)
	}

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("failed to list secrets: %s", resp.Status)
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var data map[string]interface{}
	if err := json.Unmarshal(body, &data); err != nil {
		return nil, err
	}

	keys := data["data"].(map[string]interface{})["keys"].([]interface{})
	secrets := make([]string, len(keys))
	for i, key := range keys {
		secrets[i] = key.(string)
	}

	return secrets, nil
}

func filterSecrets(secrets []string, query string) []string {
	var filtered []string
	for _, secret := range secrets {
		if strings.Contains(secret, query) {
			filtered = append(filtered, secret)
		}
	}
	return filtered
}

func getSecret(vaultAddr, secretEngine, secretPath, token, namespace string) (string, error) {
	url := fmt.Sprintf("%s/v1/%s/data/%s", vaultAddr, secretEngine, secretPath)
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return "", err
	}

	req.Header.Add("X-Vault-Token", token)
	if namespace != "" {
		req.Header.Add("X-Vault-Namespace", namespace)
	}

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("failed to get secret: %s", resp.Status)
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	var data map[string]interface{}
	if err := json.Unmarshal(body, &data); err != nil {
		return "", err
	}

	secretData := data["data"].(map[string]interface{})["data"].(map[string]interface{})
	secret, err := json.MarshalIndent(secretData, "", "  ")
	if err != nil {
		return "", err
	}

	return string(secret), nil
}

func displaySecret(secret, format string) {
	switch format {
	case "json":
		fmt.Println(secret)
	case "text":
		var data map[string]interface{}
		if err := json.Unmarshal([]byte(secret), &data); err != nil {
			log.Printf("Error parsing secret: %v", err)
			return
		}
		for key, value := range data {
			fmt.Printf("%s: %v\n", key, value)
		}
	default:
		fmt.Println(secret)
	}
}

func addSecret(vaultAddr, secretEngine, token, namespace string) {
	keyPrompt := promptui.Prompt{
		Label: "Secret Key",
	}

	key, err := keyPrompt.Run()
	if err != nil {
		log.Printf("Prompt failed %v\n", err)
		return
	}

	valuePrompt := promptui.Prompt{
		Label: "Secret Value",
	}

	value, err := valuePrompt.Run()
	if err != nil {
		log.Printf("Prompt failed %v\n", err)
		return
	}

	secretData := map[string]interface{}{
		key: value,
	}

	err = writeSecret(vaultAddr, secretEngine, key, secretData, token, namespace)
	if err != nil {
		log.Printf("Error adding secret: %v", err)
	} else {
		fmt.Println("Secret added successfully")
	}
}

func updateSecret(vaultAddr, secretEngine, token, namespace string) {
	keyPrompt := promptui.Prompt{
		Label: "Secret Key",
	}

	key, err := keyPrompt.Run()
	if err != nil {
		log.Printf("Prompt failed %v\n", err)
		return
	}

	valuePrompt := promptui.Prompt{
		Label: "New Secret Value",
	}

	value, err := valuePrompt.Run()
	if err != nil {
		log.Printf("Prompt failed %v\n", err)
		return
	}

	secretData := map[string]interface{}{
		key: value,
	}

	err = writeSecret(vaultAddr, secretEngine, key, secretData, token, namespace)
	if err != nil {
		log.Printf("Error updating secret: %v", err)
	} else {
		fmt.Println("Secret updated successfully")
	}
}

func writeSecret(vaultAddr, secretEngine, secretPath string, secretData map[string]interface{}, token, namespace string) error {
	url := fmt.Sprintf("%s/v1/%s/data/%s", vaultAddr, secretEngine, secretPath)
	data := map[string]interface{}{
		"data": secretData,
	}
	jsonData, err := json.Marshal(data)
	if err != nil {
		return err
	}

	req, err := http.NewRequest("POST", url, strings.NewReader(string(jsonData)))
	if err != nil {
		return err
	}

	req.Header.Add("X-Vault-Token", token)
	if namespace != "" {
		req.Header.Add("X-Vault-Namespace", namespace)
	}
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("failed to write secret: %s", resp.Status)
	}

	return nil
}
