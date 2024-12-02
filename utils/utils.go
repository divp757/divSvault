package utils

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
)

func VaultRequest(url string, mathod string, token string, namespace string, payload *strings.Reader) (map[string]interface{}, error) {
	client := &http.Client{}
	var req *http.Request
	var err error

	if mathod != "POST" {
		req, err = http.NewRequest(mathod, url, nil)

	} else {
		req, err = http.NewRequest(mathod, url, payload)
	}

	if err != nil {
		return nil, err
	}

	req.Header.Add("X-Vault-Token", token)
	req.Header.Set("Content-Type", "application/json")

	if namespace != "" {
		req.Header.Add("X-Vault-Namespace", namespace)
	}

	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("failed to %s secrets: %s", mathod, resp.Status)
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var data map[string]interface{}
	if err := json.Unmarshal(body, &data); err != nil {
		return nil, err
	}
	return data, nil
}
