package httpmanager

import (
	"fmt"
	"io"
	"log"
	"net/http"

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
	body, err := io.ReadAll(r.Body)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		_, err = w.Write([]byte("Internal server error"))
		if err != nil {
			log.Println("Error while sending error response:", err)
		}
	}

	fmt.Printf("Request to anaylze: %s\n", body)

	//w.Header().Set("Content-Type", "text/html")
	w.Write([]byte("0"))
}

func (s *Service) Run() {
	fmt.Println("Listening on", s.listenURL)
	log.Fatal(http.ListenAndServe(s.listenURL, s.mux))
}
