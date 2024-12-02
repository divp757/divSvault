package main

import (
	"divSvault/utils"
	"divSvault/vault"
	"flag"
	"log"
	"os"

	"github.com/manifoldco/promptui"
)

func getEnv(key, fallback string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return fallback
}

func main() {
	// Define flags for the CLI
	vaultAddr := flag.String("vault-addr", getEnv("VAULT_ADDR", "http://127.0.0.1:8200"), "Vault server address")
	secretEngine := flag.String("secret-engine", "", "Secret engine in Vault")
	token := flag.String("token", getEnv("VAULT_TOKEN", ""), "Vault token")
	namespace := flag.String("namespace", getEnv("VAULT_NAMESPACE", ""), "Vault namespace")
	outputFormat := flag.String("output-format", "json", "Output format (json or text)")
	apiPath := ""

	flag.Parse()

	config, err := utils.LoadConfig("configs/config.json")

	if err != nil {
		log.Fatalf("Error loading config: %v", err)
	}
	vaultConfig, err := vault.SelectVault(config.Vaults)
	if err != nil {
		log.Fatalf("Error selecting vault: %v", err)
	}

	if *token == "" {
		*token, err = vault.LdapLogin(vaultConfig)
		if err != nil {
			log.Fatalf("Error logging in: %v", err)
		}
	}

	if *secretEngine == "" {
		*secretEngine = vaultConfig.SecretEngine
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
			vault.ListAndDisplaySecrets(*vaultAddr, *secretEngine, *token, apiPath, *namespace, *outputFormat)
		case "Add Secret":
			vault.AddSecret(*vaultAddr, *secretEngine, *token, *namespace)
		case "Update Secret":
			vault.UpdateSecret(*vaultAddr, *secretEngine, *token, *namespace)
		}
	}
}
