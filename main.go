package main

import (
	"log"
	"os"

	"github.com/BlockscapeLab/goz-data-backend/chain"
	"github.com/BlockscapeLab/goz-data-backend/data"
	"github.com/BlockscapeLab/goz-data-backend/rest"
)

func main() {
	lcd, err := chain.CreateLCDConnector("127.0.0.1", 1317)
	if err != nil {
		log.Println("Error while setting up connection to lcd:", err)
		os.Exit(1)
	}

	dp, err := data.NewDataProvider("2020-05-14T00:00:00Z", "2020-05-28T00:00:00Z", lcd)
	if err != nil {
		log.Println("Couldn't set up DataProvider:", err)
		os.Exit(1)
	}
	//dp := data.MockProvider{}
	rest.StartRestServer("0.0.0.0", 8080, dp)
}
