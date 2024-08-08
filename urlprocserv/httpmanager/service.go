package httpmanager

import (
	"bufio"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"strconv"

	"github.com/lstep/rspamd-url-behaviour-extension/urlprocserv/internal/urlmanager"
)

type Service struct {
	listenURL  string
	urlService *urlmanager.Service

	mux *http.ServeMux
}

func New(listenURL string, urlService *urlmanager.Service) *Service {
	return &Service{
		listenURL:  listenURL,
		urlService: urlService,
	}
}

func (s *Service) SetupRoutes() {
	s.mux = http.NewServeMux()
	s.mux.HandleFunc("/", s.CheckURLs)
}

func (s *Service) CheckURLs(w http.ResponseWriter, r *http.Request) {
	dataScanner := bufio.NewScanner(r.Body)
	dataScanner.Split(bufio.ScanLines)
	defer r.Body.Close()

	var ipAddress string
	var email string
	var urlList []string

	for dataScanner.Scan() {
		if ipAddress == "" {
			ipAddress = dataScanner.Text()
			continue
		} else if email == "" {
			email = dataScanner.Text()
			continue
		}

		url := dataScanner.Text()
		if url != "" && IsUrl(url) {
			urlList = append(urlList, url)
		}
	}

	fmt.Printf("Request to anaylze: %#v\n", urlList)
	result := strconv.Itoa(s.urlService.CheckURLs(urlList))

	//w.Header().Set("Content-Type", "text/html")
	w.Write([]byte(result))
}

func (s *Service) Run() {
	fmt.Println("Listening on", s.listenURL)
	log.Fatal(http.ListenAndServe(s.listenURL, s.mux))
}

// IsUrl checks if the input string is a valid URL.
//
// str: input string to validate.
// bool: returns true if the input is a valid URL, false otherwise.
func IsUrl(str string) bool {
	u, err := url.Parse(str)
	return err == nil && u.Scheme != "" && u.Host != ""
}
