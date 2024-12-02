package vault

import (
	"fmt"
	"log"

	"github.com/manifoldco/promptui"
)

func UpdateSecret(vaultAddr, secretEngine, token, namespace string) {
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

	err = WriteSecret(vaultAddr, secretEngine, key, secretData, token, namespace)
	if err != nil {
		log.Printf("Error updating secret: %v", err)
	} else {
		fmt.Println("Secret updated successfully")
	}
}
