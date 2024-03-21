package main

import "flag"

var (
	listenURL *string
)

func main() {
	listenURL = flag.String("listen", "127.0.0.1:8088", "listen url")

	flag.Parse()

	urlService := urlmanager.New()

	server := httpmanager.New(*listenURL, urlService)
	server.SetupRoutes()
	server.Run()
}
