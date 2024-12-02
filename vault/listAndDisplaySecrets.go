package vault

import (
	"log"
	"strings"

	"github.com/manifoldco/promptui"
)

func ListAndDisplaySecrets(vaultAddr, secretEngine, token, apiPath, namespace, outputFormat string) {
	secrets, err := ListSecrets(vaultAddr, secretEngine, token, apiPath, namespace)
	if err != nil {
		log.Printf("Error listing secrets: %v", err)
		return
	}

	searchPrompt := promptui.Prompt{
		Label: "Search Secret or press Enter to see all the secres. > " + secretEngine + "/" + apiPath,
	}

	searchQuery, err := searchPrompt.Run()
	if err != nil {
		log.Printf("Prompt failed %v\n", err)
		return
	}

	filteredSecrets := FilterSecrets(secrets, searchQuery)

	selectPrompt := promptui.Select{
		Label: "Select Secret",
		Items: filteredSecrets,
	}

	_, result, err := selectPrompt.Run()
	if err != nil {
		log.Printf("Prompt failed %v\n", err)
		return
	}

	apiPath = apiPath + result

	if strings.Contains(result, "/") {
		ListAndDisplaySecrets(vaultAddr, secretEngine, token, apiPath, namespace, outputFormat)
	} else {
		secret, err := GetSecret(vaultAddr, secretEngine, apiPath, token, namespace)
		if err != nil {
			log.Printf("Error getting secret: %v", err)
			return
		}
		DisplaySecret(secret, outputFormat)
	}
}
