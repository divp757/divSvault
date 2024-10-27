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

	// // Set up signal handling for graceful exit
	// sigs := make(chan os.Signal, 1)
	// signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)

	// go func() {
	// 	sig := <-sigs
	// 	fmt.Printf("\nReceived signal: %s. Exiting...\n", sig)
	// 	os.Exit(0)
	// }()

	for {
		// Fetch the list of secrets
		secrets, err := listSecrets(*vaultAddr, *secretEngine, *token, *namespace)
		if err != nil {
			log.Printf("Error listing secrets: %v", err)
			continue
		}

		// Create a search prompt to filter secrets
		searchPrompt := promptui.Prompt{
			Label: "Search Secret",
		}

		searchQuery, err := searchPrompt.Run()
		if err != nil {
			log.Printf("Prompt failed %v\n", err)
			os.Exit(0)
		}

		filteredSecrets := filterSecrets(secrets, searchQuery)

		// Create a prompt to select a secret
		selectPrompt := promptui.Select{
			Label: "Select Secret",
			Items: filteredSecrets,
		}

		_, result, err := selectPrompt.Run()
		if err != nil {
			log.Printf("Prompt failed %v\n", err)
			os.Exit(0)
		}

		// Fetch and display the selected secret
		secret, err := getSecret(*vaultAddr, *secretEngine, result, *token, *namespace)
		if err != nil {
			log.Printf("Error getting secret: %v", err)
			continue
		}

		displaySecret(secret, *outputFormat)
	}
}

func getEnv(key, fallback string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return fallback
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
