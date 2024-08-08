package urlmanager

import (
	"errors"
	"fmt"
	"net/http"
	"strings"
	"time"
)

type Service struct {
}

func New() *Service {
	return &Service{}
}

func (s *Service) CheckURLs(urls []string) int {
	for _, url := range urls {
		s.CheckURL(url)
	}

	// @TODO: calculate total
	return 0
}

func (s *Service) CheckURL(url string) int {
	nbRedirects, err := s.getNbRedirects(url)
	if err != nil {
		fmt.Println(err)
		return 0
	}

	return nbRedirects
}

func (s *Service) getNbRedirects(url string) (int, error) {
	var nbRedirects int

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return 0, err
	}

	client := new(http.Client)
	client.CheckRedirect = func(req *http.Request, via []*http.Request) error {
		nbRedirects++
		return errors.New("Redirect")
	}

	response, err := client.Do(req)
	if response != nil && response.StatusCode == http.StatusFound { //status code 302
		fmt.Println(response.Location())
	} else { // handle err or response == nil
		fmt.Printf("rsp: %+v, err: %v\n", response, err)
	}
	response.Body.Close()

	return nbRedirects, nil
}

func (s *Service) getFinalDestination(url string) (string, error) {
	client := &http.Client{
		Timeout: 15 * time.Second,
	}

	resp, err := client.Get(url)
	if err != nil {
		return "", err
	}

	finalURL := resp.Request.URL.String()
	resp.Body.Close()

	return finalURL, nil
}

// Need to check the Finaldestination, not the first one
func (s *Service) isReachable(url string) bool {
	timeout := time.Duration(2 * time.Second)
	client := http.Client{
		Timeout: timeout,
	}
	resp, err := client.Get(url)
	if err != nil {
		return false
	}
	defer resp.Body.Close()

	return resp.StatusCode == 200
}

func (s *Service) getSSLInfos(url string) (string, error) {
	if !strings.HasPrefix(url, "https://") {
		return "", errors.New("URL must start with https://")
	}

	resp, err := http.Get(url)
	if err != nil {
		return "", errors.New("Failed to get URL")
	}
	defer resp.Body.Close()

	return resp.TLS.PeerCertificates[0].Issuer.String(), nil
}
