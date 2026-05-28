package webclient

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"time"

	"github.com/sverrehu/versions/internal/config"
)

var client = &http.Client{
	Timeout: 60 * time.Second,
}

type HTTPError struct {
	StatusCode int
	Status     string
	Body       string
}

func (e *HTTPError) Error() string {
	if e.Body == "" {
		return e.Status
	}
	return fmt.Sprintf("%s: %s", e.Status, e.Body)
}

func Get(url string, credentials *config.Credentials) (string, error) {
	log.Printf("fetching %s", url)
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return "", err
	}
	addCredentials(req, credentials)
	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}
	if resp.StatusCode == http.StatusNotFound {
		return "", nil
	}
	if resp.StatusCode != http.StatusOK {
		return "", &HTTPError{
			StatusCode: resp.StatusCode,
			Status:     resp.Status,
			Body:       string(body),
		}
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
