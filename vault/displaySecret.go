package vault

import (
	"encoding/json"
	"fmt"
	"log"
)

func DisplaySecret(secret, format string) {
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
