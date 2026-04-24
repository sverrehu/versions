package webclient

import (
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"

	"github.com/sverrehu/versions/internal/config"
)

func Get(url string, credentials *config.Credentials) (string, error) {
	log.Printf("fetching %s", url)
	client := &http.Client{}
	req, _ := http.NewRequest("GET", url, nil)
	addCredentials(req, credentials)
	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()
	if resp.StatusCode == http.StatusNotFound {
		return "", nil
	}
	if resp.StatusCode != http.StatusOK {
		return "", errors.New(resp.Status)
	}
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}
	return string(body), nil
}

func addCredentials(req *http.Request, credentials *config.Credentials) {
	if credentials == nil {
		return
	}
	if credentials.Token != "" {
		req.Header.Add("Authorization", fmt.Sprintf("Bearer %s", credentials.Token))
	}
}
