package rest

import (
	"fmt"
	"log"
	"net/http"
)

//StartRestServer start the server at specified port and ip.
func StartRestServer(ip string, port int) {
	registerHandlers()

	go startServer(ip, port)
}

func startServer(ip string, port int) {
	if err := http.ListenAndServe(fmt.Sprintf("http://%s:%d", ip, port), nil); err != nil {
		log.Println("[Error] http server failed:", err)
	}

}

func registerHandlers() {
	// http.HandleFunc("/test", testHandler)
}
