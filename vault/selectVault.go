package vault

import (
	"divSvault/utils"
	"fmt"

	"github.com/manifoldco/promptui"
)

func SelectVault(vaults []utils.VaultConfig) (utils.VaultConfig, error) {
	vaultNames := make([]string, len(vaults))
	for i, v := range vaults {
		vaultNames[i] = v.Name
	}

	vaultPrompt := promptui.Select{
		Label: "Select Vault",
		Items: vaultNames,
	}

	_, result, err := vaultPrompt.Run()
	if err != nil {
		return utils.VaultConfig{}, err
	}

	for _, v := range vaults {
		if v.Name == result {
			return v, nil
		}
	}

	return utils.VaultConfig{}, fmt.Errorf("vault not found")
}
