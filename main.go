package main

import (
	"github.com/BlockscapeLab/goz-data-backend/data"
	"github.com/BlockscapeLab/goz-data-backend/rest"
)

func main() {
	//dp := &data.DataProvider{}
	dp := data.MockProvider{}
	rest.StartRestServer("0.0.0.0", 8080, dp)
}
